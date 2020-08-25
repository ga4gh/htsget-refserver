// Package htsconstants contains program constants
//
// Module endpoints_test tests module endpoints
package htsconstants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var endpointsTC = []struct {
	e   APIEndpoint
	exp string
}{
	{APIEndpointReadsTicket, "/reads/{id}"},
	{APIEndpointReadsData, "/reads/data/{id}"},
	{APIEndpointVariantsServiceInfo, "/variants/service-info"},
	{APIEndpointFileBytes, "/file-bytes"},
}

func TestEndpoints(t *testing.T) {
	for _, tc := range endpointsTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}
