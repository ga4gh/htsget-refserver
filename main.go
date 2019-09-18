package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)

var dataSource = "http://s3.amazonaws.com/czbiohub-tabula-muris/"

type Ticket struct {
	HTSget Container `json:"htsget"`
}

type Container struct {
	Format string `json:"format"`
	URLS   []URL  `json:"urls"`
}

type URL struct {
	URL     string  `json:"url"`
	Headers Headers `json:"headers"`
	Class   string  `json:"class"`
}

type Headers struct {
	Range string `json:"range"`
}

func main() {
	router := chi.NewRouter()

	// serve index.html at root of api
	staticPath, _ := filepath.Abs("./")
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/", fs)

	// route for "reads" resource
	router.Get("/reads/{id}", getReads)

	http.ListenAndServe(":3000", router)
}

func getReads(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	// *** Parse query params ***
	params := req.URL.Query()

	// format param - optional
	var format string
	if _, ok := params["format"]; ok {
		if validReadFormat(params["format"][0]) {
			format = strings.ToUpper(params["format"][0])
		} else {
			panic("UnsupportedFormat")
		}
	} else {
		format = "BAM"
	}

	// class param
	var class string
	if _, ok := params["class"]; ok {
		if validClass(params["class"][0]) {
			class = strings.ToLower(params["class"][0])
		} else {
			panic("InvalidInput")
		}
	}

	// referenceName param
	var referenceName string
	if _, ok := params["referenceName"]; ok {
		referenceName = params["referenceName"][0]
	}

	var fileName string
	if strings.HasPrefix(id, "10X") {
		fileName = "10x_bam_files/" + id
	} else {
		fileName = "facs_bam_files/" + id
	}
	urls := []URL{{dataSource + fileName, Headers{"bytes=1-100"}, class}}
	container := Container{format, urls}
	ticket := Ticket{HTSget: container}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.ga4gh.htsget.v1.0.0+json; charset=utf-8")
	w.Write(ticketJSON)
}

func validReadFormat(s string) bool {
	switch strings.ToUpper(s) {
	case "BAM":
		return true
	case "CRAM":
		return true
	default:
		return false
	}
}

func validClass(s string) bool {
	switch strings.ToLower(s) {
	case "head":
		return true
	case "body":
		return true
	default:
		return false
	}
}
