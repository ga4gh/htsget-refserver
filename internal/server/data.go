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
	_, err := parseFormat(params) // server currently only supports BAM
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"UnsupportedFormat",
				"The requested file format is not supported by the server",
			},
		}
		writeError(w, htsErr)
		return
	}

	refName := parseRefName(params)
	start, end, err := parseRange(params, refName)
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"InvalidRange",
				"The request range cannot be satisfied",
			},
		}
		writeError(w, htsErr)
		return
	}

	fields, err := parseFields(params)
	if !strings.HasPrefix(id, "10X") {
		fields, err = parseFields(params)
		if err != nil {
			htsErr := &htsgetError{
				Code: http.StatusBadRequest,
				Htsget: errorContainer{
					"InvalidInput",
					"The request parameters do not adhere to the specification",
				},
			}
			writeError(w, htsErr)
			return
		}
	}

	blockID, err := strconv.Atoi(r.Header.Get("block-id"))
	if err != nil {
		writeError(w, err)
		return
	}
	numBlocks, err := strconv.Atoi(r.Header.Get("num-blocks"))
	if err != nil {
		writeError(w, err)
		return
	}

	class := r.Header.Get("class")
	region := &genomics.Region{Name: refName, Start: start, End: end}

	args := getCmdArgs(id, region, class, fields)
	cmd := exec.Command("samtools", args...)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		writeError(w, err)
		return
	}

	err = cmd.Start()
	if err != nil {
		writeError(w, err)
		return
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
			headerLen, err := headerLen(id)
			w.Header().Set("header-len", strconv.FormatInt(headerLen, 10))
			if err != nil {
				writeError(w, err)
				return
			}
			headerBuf := make([]byte, headerLen)
			io.ReadFull(reader, headerBuf)
		}
		if blockID != numBlocks { // remove EOF if current block is not the last block
			bufSize := 65536
			buf := make([]byte, bufSize)
			n, err := io.ReadFull(reader, buf)
			if err != nil && err.Error() != "unexpected EOF" {
				writeError(w, err)
				return
			}

			eofBuf := make([]byte, eofLen)
			for n == bufSize {
				copy(eofBuf, buf[n-eofLen:])
				w.Write(buf[:n-eofLen])
				n, err = io.ReadFull(reader, buf)
				if err != nil && err.Error() != "unexpected EOF" {
					writeError(w, err)
					return
				}
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

		tmpDirPath, err := tmpDirPath()
		if err != nil {
			writeError(w, err)
			return
		}
		tmpPath := tmpDirPath + id
		tmp, err := os.Create(tmpPath)
		if err != nil {
			writeError(w, err)
			return
		}

		l, _, eof := reader.ReadLine()
		for ; eof == nil; l, _, eof = reader.ReadLine() {
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
			} else {
				l = append(l, "\n"...)
			}
			_, err = tmp.Write(l)
			if err != nil {
				writeError(w, err)
				return
			}
		}
		tmp.Close()
		bamCmd := exec.Command("samtools", "view", "-b", tmpPath)
		bamPipe, err := bamCmd.StdoutPipe()
		if err != nil {
			writeError(w, err)
			return
		}
		err = bamCmd.Start()
		if err != nil {
			writeError(w, err)
			return
		}

		emptyHeaderLen := 50
		// remove header
		bamReader := bufio.NewReader(bamPipe)
		headerBuf := make([]byte, emptyHeaderLen)
		io.ReadFull(bamReader, headerBuf)
		io.Copy(w, bamReader)

		err = bamCmd.Wait()
		if err != nil {
			writeError(w, err)
			return
		}

		err = os.Remove(tmpPath)
		if err != nil {
			writeError(w, err)
			return
		}
	}
	cmd.Wait()
}

func getTempPath(id string, blockID int) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(cwd)
	tempPath := parent + "/temp/" + id + "_" + strconv.Itoa(blockID)

	exists, err := fileExists(parent + "/temp")
	if err != nil {
		return "", err
	}
	if !exists {
		os.Mkdir(parent+"/temp/", 0755)
	} else {
		isDir, err := isDir(parent + "/temp")
		if err != nil {
			return "", err
		}
		if !isDir {
			os.Mkdir(parent+"/temp/", 0755)
		}
	}
	return tempPath, nil
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

func headerLen(id string) (int64, error) {
	cmd := exec.Command("samtools", "view", "-H", "-b", dataSource+filePath(id))
	tmpDirPath, err := tmpDirPath()
	if err != nil {
		return 0, err
	}
	path := tmpDirPath + id + "_header"
	tmpHeader, err := os.Create(path)
	if err != nil {
		return 0, err
	}

	cmd.Stdout = tmpHeader
	cmd.Run()

	fi, err := tmpHeader.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size() - 12
	tmpHeader.Close()
	os.Remove(path)
	return size, nil
}

func removeHeader(path string) error {
	fin, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fin.Close()
	fin.Seek(0, 0)
	b, err := bam.NewReader(fin, 0)
	if err != nil {
		return err
	}

	b.Header()
	b.Read()
	lastChunk := b.LastChunk()
	hLen := lastChunk.Begin.File

	fDest, err := os.Create(path + "_copy")
	if err != nil {
		return err
	}

	fin.Seek(hLen, 0)
	io.Copy(fin, fDest)

	os.Remove(path)
	os.Rename(path+"_copy", path)
	return nil
}

func removeEOF(tempPath string) error {
	fi, err := os.Stat(tempPath)
	if err != nil {
		return err
	}
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
	if err != nil {
		return false, err
	}

	if src.Mode().IsRegular() {
		fmt.Println(path + " already exist as a file!")
		return false, nil
	}
	return true, nil
}

func tmpDirPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd + "/temp/", nil
}
