package handler

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
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
	class, err := parseClass(params)
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	blockID := r.Header.Get["block-ID"]

	args := []string{"view", dataSource + filePath(id)}
	var refRange string
	var cmd *exec.Cmd
	if refName != "" {
		refRange = refName + ":" + start + "-" + end
		args = append(args, refRange)
	}
	if class == "header" {
		args = append(args, "-H")
	}
	cmd = exec.Command("samtools", args...)

	pipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	cwd, _ := os.Getwd()
	tempPath := cwd + "/temp/" + id + "_" + blockID
	fSam, _ := os.Create(tempPath)
	defer fSam.Close()
	reader := bufio.NewReader(pipe)

	if len(fields) == 0 {
		io.Copy(fSam, reader)
	} else {
		l, _, err := reader.ReadLine()
		columns := make([]int, 12)
		for _, field := range fields {
			columns[FIELDS[field]] = 1
		}
		sort.Ints(columns)

		for ; err == nil; l, _, err = reader.ReadLine() {
			if l[0] == 64 {
				l = append(l, "\n"...)
				fSam.Write(l)
			} else {
				var output []string
				ls := strings.Split(string(l), "\t")
				for i, col := range columns {
					if col == 1 {
						output = append(output, ls[i-1])
					} else {
						if i == 2 || i == 4 || i == 5 || i == 8 || i == 9 {
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

	fin, _ := os.Open(tempPath)
	defer fin.Close()
	b, _ := bam.NewReader(fin, 0)
	defer b.Close()
	b.Header()
	b.Read()
	lastChunk := b.LastChunk()
	hLen := lastChunk.Begin.File

	cmd.Wait()
}

func samToBam() {
	cmd = exec.Command("samtools", "view", "-h", "-b", tempPath, "-o", tempPath)
	cmd.Run()
}
