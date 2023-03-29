package htsserver

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/assumerole"
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	log "github.com/sirupsen/logrus"
)

// SetRouter sets up and returns a go-chi router to caller
func SetRouter() (*chi.Mux, error) {
	router := chi.NewRouter()
	log.Debug("in the setRouter call")

	// Setup CORS
	corsAllowedHeaders := strings.Split(htsconfig.GetCorsAllowedHeaders(), ",")
	allowedHeaders := append(corsAllowedHeaders, "HtsgetBlockClass", "HtsgetCurrentBlock", "HtsgetTotalBlocks", "HtsgetFilePath")
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(htsconfig.GetCorsAllowedOrigins(), ","),
		AllowedMethods:   strings.Split(htsconfig.GetCorsAllowedMethods(), ","),
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: htsconfig.GetCorsAllowCredentials(),
		MaxAge:           htsconfig.GetCorsMaxAge(),
	}))

	// Setup AWS AssumeRole middleware
	if htsconfig.IsAwsAssumeRole() {
		router.Use(assumerole.Handler(assumerole.Options{
			Debug: false,
		}))
	}

	// Add API Routes

	// if reads enabled, add reads routes
	if htsconfig.IsEndpointEnabled(htsconstants.APIEndpointReadsTicket) {
		log.Debugf("is reads enabled? %v", htsconstants.APIEndpointReadsTicket.String())

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
			log.Debug("error adding the static files route, %v", err)
			return nil, err
		}
		http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(absDocsDir))))
	}

	return router, nil
}
