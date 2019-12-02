package server

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/david-xliu/htsget-refserver/internal/genomics"
	"github.com/go-chi/chi"
)

// getData serves the actual data from AWS back to client
func getData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// *** Parse query params ***
	params := r.URL.Query()
	_, err := parseFormat(params)
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	blockID, _ := strconv.Atoi(r.Header.Get("block-id"))
	numBlocks, _ := strconv.Atoi(r.Header.Get("num-blocks"))
	class := r.Header.Get("class")

	region := &genomics.Region{Name: refName, Start: start, End: end}

	args := getCmdArgs(id, region, class, fields)
	cmd := exec.Command("samtools", args...)
	fmt.Println(args)
	pipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)

	var eofLen int
	if class == "header" {
		eofLen = HEADER_EOF_LEN
	} else {
		eofLen = EOF_LEN
	}

	if len(fields) == 0 || class == "header" {
		if class != "header" { // remove header
			headerBuf := make([]byte, headerLen(id))
			io.ReadFull(reader, headerBuf)
		}

		if blockID != numBlocks { // remove EOF if current block is not the last block
			bufSize := 65536
			buf := make([]byte, bufSize)
			n, _ := io.ReadFull(reader, buf)
			eofBuf := make([]byte, eofLen)
			for n == bufSize {
				copy(eofBuf, buf[n-eofLen:])
				w.Write(buf[:n-eofLen])
				n, _ = io.ReadFull(reader, buf)
				if n == bufSize {
					w.Write(eofBuf)
				}
			}

			if n >= eofLen {
				w.Write(buf[:n-eofLen])
			} else {
				w.Write(eofBuf[:eofLen-n])
			}
		} else {
			io.Copy(w, reader)
		}
	} else {
		columns := make([]bool, 11)
		for _, field := range fields {
			columns[FIELDS[field]-1] = true
		}

		var eof bool = false
		for !eof {
			tmpPath := tmpDirPath() + id
			tmp, err := os.Create(tmpPath)
			if err != nil {
				panic(err)
			}

			/* for i := 0; i < 500; i++ {*/
			for true {
				l, _, err := reader.ReadLine()
				if err != nil {
					eof = true
					break
				}

				if l[0] != 64 {
					var output []string
					ls := strings.Split(string(l), "\t")

					for i, col := range columns {
						if col {
							output = append(output, ls[i])
						} else {
							if i == 1 || i == 3 || i == 4 || i == 7 || i == 8 {
								output = append(output, "0")
							} else {
								output = append(output, "*")
							}
						}
					}
					l = []byte(strings.Join(output, "\t") + "\n")
					if err != nil {
						panic(err)
					}
				}
				_, err = tmp.Write(l)
			}
			tmp.Close()

			fmt.Println(eofLen)
			bamCmd := exec.Command("samtools", "view", "-b", tmpPath)
			bamPipe, _ := bamCmd.StdoutPipe()
			err = bamCmd.Start()
			if err != nil {
				panic(err)
			}

			emptyHeaderLen := 50
			// remove header
			bamReader := bufio.NewReader(bamPipe)
			headerBuf := make([]byte, emptyHeaderLen)
			io.ReadFull(bamReader, headerBuf)

			// remove EOF
			cwd, _ := os.Getwd()
			eofPath := cwd + "/eof_" + id
			eofRemoval, _ := os.Create(eofPath)

			io.Copy(eofRemoval, bamReader)
			eofRemoval.Close()
			removeEOF(eofPath)
			eofRemoval, _ = os.Open(eofPath)
			io.Copy(w, eofRemoval)

			/*     bufSize := 65536*/
			//buf := make([]byte, bufSize)
			//n, _ := io.ReadFull(bamReader, buf)
			//eofBuf := make([]byte, eofLen)
			//for n == bufSize {
			//copy(eofBuf, buf[n-eofLen:])
			//w.Write(buf[:n-eofLen])
			//n, _ = io.ReadFull(bamReader, buf)
			//if n == bufSize {
			//w.Write(eofBuf)
			//}
			//}

			//if n >= eofLen {
			//w.Write(buf[:n-eofLen])
			//} else {
			//w.Write(eofBuf[:eofLen-n])
			/*}*/

			/*     w.Write(EOF)*/
			//w.Write(EOF)
			/*w.Write(EOF)*/

			err = bamCmd.Wait()
			if err != nil {
				panic(err)
			}
			err = os.Remove(tmpPath)
			if err != nil {
				panic(err)
			}
			err = os.Remove(eofPath)
			if err != nil {
				panic(err)
			}
		}

		if blockID == numBlocks {
			w.Write(EOF)
		}
	}
	fmt.Println("exiting")
	cmd.Wait()
}

func getTempPath(id string, blockID int) string {
	cwd, _ := os.Getwd()
	parent := filepath.Dir(cwd)
	tempPath := parent + "/temp/" + id + "_" + strconv.Itoa(blockID)

	if exists, _ := fileExists(parent + "/temp"); !exists {
		os.Mkdir(parent+"/temp/", 0755)
	} else {
		if isDir, _ := isDir(parent + "/temp"); !isDir {
			os.Mkdir(parent+"/temp/", 0755)
		}
	}
	return tempPath
}

func getCmdArgs(id string, r *genomics.Region, class string, fields []string) []string {
	args := []string{"view", dataSource + filePath(id)}
	if class == "header" {
		args = append(args, "-H")
		args = append(args, "-b")
	} else {
		if len(fields) == 0 {
			args = append(args, "-b")
		}
		if r.String() != "" {
			args = append(args, r.String())
		}
	}
	return args
}

func samToBam(tempPath string) string {
	bamPath := tempPath + "_bam"
	cmd := exec.Command("samtools", "view", "-h", "-b", tempPath, "-o", bamPath)
	cmd.Run()
	return bamPath
}

func headerLen(id string) int64 {
	cmd := exec.Command("samtools", "view", "-H", "-b", dataSource+filePath(id))
	path := tmpDirPath() + id + "_header"
	tmpHeader, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	cmd.Stdout = tmpHeader
	cmd.Run()

	fi, err := tmpHeader.Stat()
	if err != nil {
		panic(err)
	}

	size := fi.Size() - 12
	tmpHeader.Close()
	os.Remove(path)
	return size
}

func removeHeader(path string) {
	fin, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	fin.Seek(0, 0)
	b, err := bam.NewReader(fin, 0)
	if err != nil {
		panic(err)
	}

	b.Header()
	b.Read()
	lastChunk := b.LastChunk()
	hLen := lastChunk.Begin.File

	fDest, err := os.Create(path + "_copy")
	if err != nil {
		panic(err)
	}

	fin.Seek(hLen, 0)
	io.Copy(fin, fDest)

	os.Remove(path)
	os.Rename(path+"_copy", path)
}

func removeEOF(tempPath string) error {
	fi, _ := os.Stat(tempPath)
	return os.Truncate(tempPath, fi.Size()-28)
}

// exists returns whether the given file or directory existsi.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}

// isDir returns whether the given path is directory
func isDir(path string) (bool, error) {
	src, err := os.Stat(path)

	if src.Mode().IsRegular() {
		fmt.Println(path + " already exist as a file!")
		return false, err
	}
	return true, err
}

func tmpDirPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	return parent + "/temp/"
}
