package main

import (
	"encoding/json"
	"net/http"
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
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the HTSget reference server :D This server serves the publicly available Tabula-Muris dataset on AWS."))
	})

	// route for "reads" resource
	r.Get("/reads/{id}", getReads)

	http.ListenAndServe(":3000", r)
}

func getReads(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var fileName string
	if strings.HasPrefix(id, "10X") {
		fileName = "10x_bam_files/" + id
	} else {
		fileName = "facs_bam_files/" + id
	}
	urls := []URL{{dataSource + fileName, Headers{"bytes=1-100"}, "body"}}
	container := Container{"BAM", urls}
	ticket := Ticket{HTSget: container}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.ga4gh.htsget.v1.0.0+json; charset=utf-8")
	w.Write(ticketJSON)
}
