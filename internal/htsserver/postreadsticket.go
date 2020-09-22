package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func postReadsTicket(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.PostMethod,
		htsconstants.APIEndpointReadsTicket,
		ticketRequestHandler,
	).handleRequest(writer, request)
}
