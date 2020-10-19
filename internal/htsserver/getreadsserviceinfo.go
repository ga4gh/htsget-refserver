package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func getReadsServiceInfo(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointReadsServiceInfo,
		noAfterSetup,
		serviceInfoRequestHandler,
	).handleRequest(writer, request)
}
