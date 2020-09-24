package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func getVariantsTicket(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointVariantsTicket,
		addRegionFromQueryString,
		ticketRequestHandler,
	).handleRequest(writer, request)
}
