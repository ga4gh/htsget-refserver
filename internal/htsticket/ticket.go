// Package htsticket produces the htsget JSON response ticket
//
// Module ticket is the base object of the htsget JSON ticket response
package htsticket

// Ticket holds the entire json ticket returned to the client
type Ticket struct {
	HTSget *Container `json:"htsget"`
}

// NewTicket instantiates an empty ticket object
func NewTicket() *Ticket {
	return new(Ticket)
}

// SetContainer sets the high-level JSON container
func (ticket *Ticket) SetContainer(container *Container) *Ticket {
	ticket.HTSget = container
	return ticket
}
