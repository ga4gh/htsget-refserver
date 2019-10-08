package handler

import (
	"bufio"
	"io"
	"net/http"
	"os/exec"
	"sort"
	"strings"

	"github.com/go-chi/chi"
)

var dataSource = "http://s3.amazonaws.com/czbiohub-tabula-muris/"

// getData serves the actual data from AWS back to client
func getData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	filePath := filePath(id)

	// *** Parse query params ***
	params := r.URL.Query()
	//format, err := parseFormat(params)
	//class, err := parseClass(params)
	//refName, err := parseRefName(params)
	//start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	fields = []string{"SEQ", "QNAME", "CIGAR"}

	// if no params are given, then directly fetch file from s3
	if len(params) == 0 {
		res, err := http.Get(dataSource + filePath)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		io.Copy(w, res.Body)
	}

	testFile := "facs_bam_files/A1-B001176-3_56_F-1-1_R1.mus.Aligned.out.sorted.bam"

	cmd := exec.Command("samtools", "view", dataSource+testFile, "chr10:10000000-12000000")
	/* cmd := exec.Command("./test.sh")*/
	/*cmd.Dir = "/Users/dliu"*/
	pipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()
	var columns []int
	for _, field := range fields {
		columns = append(columns, FIELDS[field])
	}
	sort.Ints(columns)

	for ; err == nil; l, _, err = reader.ReadLine() {
		if l[0] == 64 {
			w.Write(append(l, "\n"...))
		} else {
			var output []string
			ls := strings.Split(string(l), "\t")
			for _, col := range columns {
				output = append(output, ls[col-1])
			}
			w.Write([]byte(strings.Join(output, "\t") + "\n"))
		}
	}

	cmd.Wait()
}
