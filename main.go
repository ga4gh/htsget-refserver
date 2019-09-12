package main

import (
	"encoding/json"
	"net/http"

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
		w.Write([]byte("Welcome to the HTSget reference server :D"))
	})

	// routes for "reads" resource
	r.Route("/reads", func(r chi.Router) {
		r.Get("/", getReads) // GET /reads
	})

	http.ListenAndServe(":3000", r)
}

func getReads(w http.ResponseWriter, r *http.Request) {
	urls := []URL{{dataSource + "id", Headers{"bytes=1-100"}, "body"}}
	container := Container{"BAM", urls}
	ticket := Ticket{HTSget: container}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		panic(err)
	}

	w.Write(ticketJSON)
}
