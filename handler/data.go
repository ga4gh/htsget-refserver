package handler

import (
	"bufio"
	"net/http"
	"os/exec"
	"sort"
	"strings"

	"github.com/go-chi/chi"
)

// getData serves the actual data from AWS back to client
func getData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// *** Parse query params ***
	params := r.URL.Query()
	//format, err := parseFormat(params)
	//class, err := parseClass(params)
	//refName, err := parseRefName(params)
	//start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	fields = []string{"SEQ", "QNAME", "CIGAR"}

	cmd := exec.Command("samtools", "view", "-h", filePath(id), "chr10:10000000-12000000")
	/* cmd := exec.Command("./test.sh")*/
	/*cmd.Dir = "/Users/dliu"*/
	pipe, _ := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()
	columns := make([]int, 12)
	for _, field := range fields {
		columns[FIELDS[field]] = 1
	}
	sort.Ints(columns)

	for ; err == nil; l, _, err = reader.ReadLine() {
		if l[0] == 64 {
			w.Write(append(l, "\n"...))
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
			w.Write([]byte(strings.Join(output, "\t") + "\n"))
		}
	}

	cmd.Wait()
}
