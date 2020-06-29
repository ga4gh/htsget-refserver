package htsgetserver

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
)

// SetRouter sets up and returns a go-chi router to caller
func SetRouter() (*chi.Mux, error) {
	router := chi.NewRouter()

	// serve index.html at root of api
	staticPath, err := filepath.Abs("./")
	if err != nil {
		return nil, err
	}
	fs := http.FileServer(http.Dir(staticPath))
	router.Handle("/", fs)

	// Route for "reads" resource
	router.Get("/reads/{id}", getReadsTicket)
	router.Get("/reads/data/{id}", getReadsData)

	return router, err
}
