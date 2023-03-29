package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	log "github.com/sirupsen/logrus"
)

func getReadsTicket(writer http.ResponseWriter, request *http.Request) {
	log.Debug("get reads ticket call")
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointReadsTicket,
		addRegionFromQueryString,
		ticketRequestHandler,
	).handleRequest(writer, request)
}
