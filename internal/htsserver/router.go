package htsserver

import (
	"net/http"
	"path/filepath"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
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

	// if reads enabled, add reads routes
	if htsconfig.IsEndpointEnabled(htsconstants.APIEndpointReadsTicket) {
		router.Get(htsconstants.APIEndpointReadsTicket.String(), getReadsTicket)
		router.Post(htsconstants.APIEndpointReadsTicket.String(), postReadsTicket)
		router.Get(htsconstants.APIEndpointReadsData.String(), getReadsData)
		router.Get(htsconstants.APIEndpointReadsServiceInfo.String(), getReadsServiceInfo)
	}

	// if variants enabled, add variants routes
	if htsconfig.IsEndpointEnabled(htsconstants.APIEndpointVariantsTicket) {
		router.Get(htsconstants.APIEndpointVariantsTicket.String(), getVariantsTicket)
		router.Get(htsconstants.APIEndpointVariantsData.String(), getVariantsData)
		router.Get(htsconstants.APIEndpointVariantsServiceInfo.String(), getVariantsServiceInfo)
	}

	router.Get(htsconstants.APIEndpointFileBytes.String(), getFileBytes)
	return router, err
}
