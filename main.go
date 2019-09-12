package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

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
	w.Write([]byte(fmt.Sprintf("url:%s", "some url")))
}
