// Package htsticket produces the htsget JSON response ticket
//
// Module ticket_test tests ticket
package htsticket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ticketSetContainerTC = []struct {
	format string
	urls   []string
}{
	{
		"BAM",
		[]string{
			"http://htsget.ga4gh.org/reads/data/object1",
		},
	},
	{
		"VCF",
		[]string{
			"http://htsget.ga4gh.org/variants/data/object1",
			"http://localhost:3000/variants/data/1000genomes.00001",
			"http://localhost:4000/variants/data/gatktest.11111",
		},
	},
}

func TestTicketSetContainer(t *testing.T) {
	for _, tc := range ticketSetContainerTC {
		ticket := NewTicket()
		container := NewContainer()
		container.setFormat(tc.format)
		urls := []*URL{}
		for i := 0; i < len(tc.urls); i++ {
			url := NewURL()
			url.SetURL(tc.urls[i])
			urls = append(urls, url)
		}
		container.SetURLS(urls)
		ticket.SetContainer(container)

		assert.Equal(t, tc.format, ticket.HTSget.Format)

		for i := 0; i < len(tc.urls); i++ {
			assert.Equal(t, tc.urls[i], ticket.HTSget.URLS[i].URL)
		}
	}
}
