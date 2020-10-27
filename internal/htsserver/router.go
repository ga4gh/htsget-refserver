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
		router.Post(htsconstants.APIEndpointVariantsTicket.String(), postVariantsTicket)
		router.Get(htsconstants.APIEndpointVariantsData.String(), getVariantsData)
		router.Get(htsconstants.APIEndpointVariantsServiceInfo.String(), getVariantsServiceInfo)
	}

	// add the file bytes endpoint for streaming byte indices of local files
	router.Get(htsconstants.APIEndpointFileBytes.String(), getFileBytes)

	// add the static files route
	docsDir := htsconfig.GetDocsDir()
	if docsDir != "" {
		absDocsDir, err := filepath.Abs(docsDir)
		if err != nil {
			return nil, err
		}
		http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(absDocsDir))))
	}

	return router, nil
}
