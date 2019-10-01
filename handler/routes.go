package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"path/filepath"
)

func SetRouter() *Mux {
	router := chi.NewRouter()

	// serve index.html at root of api
	staticPath, _ := filepath.Abs("./")
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/", fs)

	// Route for "reads" resource
	router.Get("/reads/{id}", getReads)

	return router
}
