package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func getVariantsServiceInfo(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointVariantsServiceInfo,
		serviceInfoRequestHandler,
	).handleRequest(writer, request)
}
