package htsrequest

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

var requestConstructDataEndpointURLTC = []struct {
	endpoint                      htsconstants.APIEndpoint
	id, referenceName, start, end string
	fields, tags, notags          []string
	exp                           string
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"object0052",
		"chr1",
		"65000",
		"420000",
		defaultListParameterValues["fields"],
		defaultListParameterValues["tags"],
		defaultListParameterValues["notags"],
		"http://localhost:3000/reads/data/object0052?end=420000&referenceName=chr1&start=65000",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.00001",
		"chr22",
		"11000000",
		"45000000",
		[]string{"SEQ", "QUAL"},
		[]string{"NM", "HI"},
		defaultListParameterValues["notags"],
		"http://localhost:3000/reads/data/tabulamuris.00001?end=45000000&fields=SEQ%2CQUAL&referenceName=chr22&start=11000000&tags=NM%2CHI",
	},
}

func TestConstructDataEndpointURL(t *testing.T) {

	for _, tc := range requestConstructDataEndpointURLTC {
		request := NewHtsgetRequest()
		request.SetEndpoint(tc.endpoint)
		request.AddScalarParam("id", tc.id)
		request.AddScalarParam("referenceName", tc.referenceName)
		request.AddScalarParam("start", tc.start)
		request.AddScalarParam("end", tc.end)
		request.AddListParam("fields", tc.fields)
		request.AddListParam("tags", tc.tags)
		request.AddListParam("notags", tc.notags)
		ep, _ := request.ConstructDataEndpointURL()
		assert.Equal(t, tc.exp, ep.String())
	}
}
