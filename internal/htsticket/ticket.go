// Package htsticket produces the htsget JSON response ticket
//
// Module ticket is the base object of the htsget JSON ticket response
package htsticket

import (
	"encoding/json"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// Ticket holds the entire json ticket returned to the client
type Ticket struct {
	HTSget *Container `json:"htsget"`
}

// NewTicket instantiates an empty ticket object
func newTicket() *Ticket {
	return new(Ticket)
}

// SetContainer sets the high-level JSON container
func (ticket *Ticket) setContainer(container *Container) *Ticket {
	ticket.HTSget = container
	return ticket
}

// FinalizeTicket for /ticket endpoints, write the htsget ticket to the HTTP
// writer
func FinalizeTicket(format string, urls []*URL, writer http.ResponseWriter) {
	container := NewContainer().setFormat(format).SetURLS(urls)
	ticket := newTicket().setContainer(container)
	writer.Header().Set(htsconstants.ContentTypeHeader.String(), htsconstants.ContentTypeHeaderHtsgetJSON.String())
	json.NewEncoder(writer).Encode(ticket)
}
