package handler

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
)

// SetRouter sets up and returns a go-chi router to caller
func SetRouter() *chi.Mux {
	router := chi.NewRouter()

	// serve index.html at root of api
	staticPath, _ := filepath.Abs("./")
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/", fs)

	// Route for "reads" resource
	router.Get("/reads/{id}", getReads)

	return router
}
