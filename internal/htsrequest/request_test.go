package htsrequest

import (
	"fmt"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

func TestConstructDataEndpointURL(t *testing.T) {

	request := NewHtsgetRequest()
	request.SetEndpoint(htsconstants.APIEndpointReadsTicket)
	request.AddScalarParam("id", "object0069")
	request.AddScalarParam("referenceName", "chr1")
	request.AddScalarParam("start", "69000")
	request.AddScalarParam("end", "4200000")
	request.AddListParam("fields", defaultListParameterValues["fields"])
	request.AddListParam("tags", defaultListParameterValues["tags"])
	request.AddListParam("notags", defaultListParameterValues["notags"])
	ep, _ := request.ConstructDataEndpointURL()
	fmt.Println(ep)
	assert.Equal(t, "a", "b")
}
