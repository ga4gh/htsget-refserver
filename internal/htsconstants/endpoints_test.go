// Package htsconstants contains program constants
//
// Module endpoints_test tests module endpoints
package htsconstants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// endpointsStringTC test cases for String method
var endpointsStringTC = []struct {
	e   APIEndpoint
	exp string
}{
	{APIEndpointReadsTicket, "/reads/{id}"},
	{APIEndpointReadsData, "/reads/data/{id}"},
	{APIEndpointVariantsServiceInfo, "/variants/service-info"},
	{APIEndpointFileBytes, "/file-bytes"},
}

// endpointsDataEndpointPathTC test cases for DataEndpointPath
var endpointsDataEndpointPathTC = []struct {
	e   APIEndpoint
	exp string
}{
	{APIEndpointReadsTicket, "/reads/data/"},
	{APIEndpointVariantsTicket, "/variants/data/"},
}

// endpointsAllowedFormatsTC test cases for AllowedFormats
var endpointsAllowedFormatsTC = []struct {
	e   APIEndpoint
	exp []string
}{
	{APIEndpointReadsTicket, []string{"BAM"}},
	{APIEndpointReadsData, []string{"BAM"}},
	{APIEndpointVariantsTicket, []string{"VCF"}},
	{APIEndpointVariantsData, []string{"VCF"}},
}

// TestEndpointsString tests String function
func TestEndpointsString(t *testing.T) {
	for _, tc := range endpointsStringTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}

// TestEndpointsDataEndpointPath tests DataEndpointPath function
func TestEndpointsDataEndpointPath(t *testing.T) {
	for _, tc := range endpointsDataEndpointPathTC {
		assert.Equal(t, tc.exp, tc.e.DataEndpointPath())
	}
}

// TestEndpointsAllowedFormats tests AllowedFormats function
func TestEndpointsAllowedFormats(t *testing.T) {
	for _, tc := range endpointsAllowedFormatsTC {
		allowedFormats := tc.e.AllowedFormats()
		for i := 0; i < len(tc.exp); i++ {
			assert.Equal(t, tc.exp[i], allowedFormats[i])
		}
	}
}
