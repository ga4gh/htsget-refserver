// Package htsconstants contains program constants
//
// Module headers_test tests module headers
package htsconstants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// headerNameTC test cases for HeaderName String function
var headerNameTC = []struct {
	e   HTTPHeaderName
	exp string
}{
	{ContentTypeHeader, "Content-Type"},
}

// contentTypeTC test cases for ContentTypeHeader String values function
var contentTypeTC = []struct {
	e   ContentTypeHeaderValue
	exp string
}{
	{ContentTypeHeaderHtsgetJSON, "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"},
}

// TestHeaderName tests HeaderName String function
func TestHeaderName(t *testing.T) {
	for _, tc := range headerNameTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}

// TestContentTypeHeaderValues tests ContentTypeHeader String function
func TestContentTypeHeaderValues(t *testing.T) {
	for _, tc := range contentTypeTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}
