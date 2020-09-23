package htsrequest

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

func TestSetParam(t *testing.T) {
	method := htsconstants.GetMethod
	endpoint := htsconstants.APIEndpointReadsTicket
	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://htsget.org/object1?format=bam&fields=SEQ,TLEN,QUAL&tags=PU,BB,SS", nil)

	htsgetReq, _ := SetAllParameters(method, endpoint, writer, request)

	fmt.Println("***")
	fmt.Println(htsgetReq.GetFormat())
	fmt.Println(htsgetReq.GetClass())
	fmt.Println(htsgetReq.GetFields())
	fmt.Println(len(htsgetReq.GetFields()))
	fmt.Println(htsgetReq.GetTags())
	fmt.Println(len(htsgetReq.GetTags()))
	fmt.Println("***")

	assert.Equal(t, "a", "b")
}
