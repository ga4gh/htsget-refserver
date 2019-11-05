package handler

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/go-chi/chi"
)

// getData serves the actual data from AWS back to client
func getData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// *** Parse query params ***
	params := r.URL.Query()
	format, err := parseFormat(params)
	if format != "BAM" {
		panic("format not supported")
	}
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	blockID, _ := strconv.Atoi(r.Header.Get("block-id"))
	numBlocks, _ := strconv.Atoi(r.Header.Get("num-blocks"))
	class := r.Header.Get("class")

	args := []string{"view", dataSource + filePath(id)}
	var refRange string
	var cmd *exec.Cmd
	if refName != "" {
		refRange = refName + ":" + start + "-" + end
		args = append(args, refRange)
	}
	if class == "header" {
		args = append(args, "-H")
	} else {
		args = append(args, "-h")
	}
	cmd = exec.Command("samtools", args...)

	pipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	cwd, _ := os.Getwd()
	tempPath := cwd + "/temp/" + id + "_" + strconv.Itoa(blockID)

	if exists, _ := exists(cwd + "/temp"); !exists {
		os.Mkdir(cwd+"/temp/", 0755)
	} else {
		if isDir, _ := isDir(cwd + "/temp"); !isDir {
			os.Mkdir(cwd+"/temp/", 0755)
		}
	}

	fSam, _ := os.Create(tempPath)
	defer fSam.Close()
	reader := bufio.NewReader(pipe)

	if len(fields) == 0 {
		io.Copy(fSam, reader)
	} else {
		l, _, err := reader.ReadLine()
		columns := make([]bool, 11)
		for _, field := range fields {
			columns[FIELDS[field]-1] = true
		}

		for ; err == nil; l, _, err = reader.ReadLine() {
			if l[0] == 64 {
				l = append(l, "\n"...)
				fSam.Write(l)
			} else {
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
				fSam.Write(l)
			}
		}
	}
	cmd.Wait()

	samToBam(tempPath)

	if class != "" && numBlocks > 1 {
		if class != "header" {
			removeHeader(tempPath)
		}
		if blockID != numBlocks {
			removeEOF(tempPath)
		}
	}

	fclient, _ := os.Open(tempPath)
	defer fclient.Close()
	io.Copy(w, fclient)
}

func samToBam(tempPath string) {
	cmd := exec.Command("samtools", "view", "-h", "-b", tempPath, "-o", tempPath)
	cmd.Run()
}

func removeHeader(tempPath string) {
	fin, _ := os.Open(tempPath)
	defer fin.Close()
	b, _ := bam.NewReader(fin, 0)
	defer b.Close()
	b.Header()
	b.Read()
	lastChunk := b.LastChunk()
	hLen := lastChunk.Begin.File

	fDest, _ := os.Create(tempPath + "_copy")
	fin.Seek(hLen, 0)
	io.Copy(fin, fDest)

	os.Remove(tempPath)
	os.Rename(tempPath+"_copy", tempPath)
}

func removeEOF(tempPath string) error {
	fi, _ := os.Stat(tempPath)
	return os.Truncate(tempPath, fi.Size()-12)
}

// exists returns whether the given file or directory existsi.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return true, err
	}
	return false, nil
}

// isDir returne whether the given path is directory
func isDir(path string) (bool, error) {
	src, err := os.Stat(path)

	if src.Mode().IsRegular() {
		fmt.Println(path + " already exist as a file!")
		return false, err
	}
	return true, err
}
