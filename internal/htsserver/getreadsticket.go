package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func getReadsTicket(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointReadsTicket,
		ticketRequestHandler,
	).handleRequest(writer, request)
}
