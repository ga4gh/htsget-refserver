package handler

import (
	"github.com/go-chi/chi"
	"io"
	"net/http"
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

	resp, err := http.Get(dataSource + filePath)
	if err != nil {
		panic(err)
	}
	//io.Copy(w, resp.Body)
}
