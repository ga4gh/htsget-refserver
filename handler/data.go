package handler

import (
	"bufio"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"os/exec"
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

	// if no params are given, then directly fetch file from s3
	if len(params) == 0 {
		resp, err := http.Get(dataSource + filePath)
		if err != nil {
			panic(err)
		}
		io.Copy(w, resp.Body)
	}

	testFile := "facs_bam_files/A1-B001176-3_56_F-1-1_R1.mus.Aligned.out.sorted.bam"

	cmd := exec.Command("samtools", "view", "-b", dataSource+testFile)
	/* cmd := exec.Command("./test.sh")*/
	/*cmd.Dir = "/Users/dliu"*/
	pipe, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(pipe)
	/* buffer := make([]byte, 12)*/
	//n, err := reader.Read(buffer)
	//for ; err == nil; n, err = reader.Read(buffer) {
	//w.Write(buffer[:n])
	//}

	io.Copy(w, reader)

	cmd.Wait()
}
