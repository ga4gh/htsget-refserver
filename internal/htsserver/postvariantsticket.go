package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func postVariantsTicket(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.PostMethod,
		htsconstants.APIEndpointVariantsTicket,
		noAfterSetup,
		ticketRequestHandler,
	).handleRequest(writer, request)
}
