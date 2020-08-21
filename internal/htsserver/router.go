package htsserver

import (
	"net/http"
	"path/filepath"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

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

	// Add API Routes

	router.Get(htsconstants.ReadsTicket.String(), getReadsTicket)
	router.Get(htsconstants.ReadsData.String(), getReadsData)
	router.Get(htsconstants.VariantsTicket.String(), getVariantsTicket)
	router.Get("/file-bytes", getFileBytes)

	return router, err
}
