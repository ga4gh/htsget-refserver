// Package htsgetticket ...
package htsgetticket

// Ticket holds the entire json ticket returned to the client
type Ticket struct {
	HTSget *Container `json:"htsget"`
}

func NewTicket() *Ticket {
	return new(Ticket)
}

func (ticket *Ticket) SetContainer(container *Container) *Ticket {
	ticket.HTSget = container
	return ticket
}
