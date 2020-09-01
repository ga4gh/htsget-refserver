// Package htsticket produces the htsget JSON response ticket
//
// Module ticket_test tests ticket
package htsticket

import (
	"net/http/httptest"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

var ticketSetContainerTC = []struct {
	format               string
	urls                 []string
	expBody              string
	expContentTypeHeader string
}{
	{
		"BAM",
		[]string{
			"http://htsget.ga4gh.org/reads/data/object1",
		},
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://htsget.ga4gh.org/reads/data/object1\"}]}}\n",
		htsconstants.ContentTypeHeaderHtsgetJSON.String(),
	},
	{
		"VCF",
		[]string{
			"http://htsget.ga4gh.org/variants/data/object1",
			"http://localhost:3000/variants/data/1000genomes.00001",
			"http://localhost:4000/variants/data/gatktest.11111",
		},
		"{\"htsget\":{\"format\":\"VCF\",\"urls\":[{\"url\":\"http://htsget.ga4gh.org/variants/data/object1\"},{\"url\":\"http://localhost:3000/variants/data/1000genomes.00001\"},{\"url\":\"http://localhost:4000/variants/data/gatktest.11111\"}]}}\n",
		htsconstants.ContentTypeHeaderHtsgetJSON.String(),
	},
}

func TestTicketFinalizeTicket(t *testing.T) {

	for _, tc := range ticketSetContainerTC {
		writer := httptest.NewRecorder()
		urls := []*URL{}
		for i := 0; i < len(tc.urls); i++ {
			url := NewURL()
			url.SetURL(tc.urls[i])
			urls = append(urls, url)
		}

		FinalizeTicket(tc.format, urls, writer)
		assert.Equal(t, tc.expBody, writer.Body.String())
		assert.Equal(t, tc.expContentTypeHeader, writer.HeaderMap[htsconstants.ContentTypeHeader.String()][0])

		// assert.Equal(t, "a", "b")

		/*
			ticket := NewTicket()
			container := NewContainer()
			container.setFormat(tc.format)

			container.SetURLS(urls)
			ticket.SetContainer(container)

			assert.Equal(t, tc.format, ticket.HTSget.Format)

			for i := 0; i < len(tc.urls); i++ {
				assert.Equal(t, tc.urls[i], ticket.HTSget.URLS[i].URL)
			}
		*/
	}
}
