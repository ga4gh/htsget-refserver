// Package htsconstants contains program constants
//
// Module headers_test tests module headers
package htsconstants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var headerNameTC = []struct {
	e   HttpHeaderName
	exp string
}{
	{ContentTypeHeader, "Content-Type"},
}

var contentTypeTC = []struct {
	e   ContentTypeHeaderValue
	exp string
}{
	{ContentTypeHeaderHtsgetJSON, "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"},
}

func TestHeaderName(t *testing.T) {
	for _, tc := range headerNameTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}

func TestContentTypeHeaderValues(t *testing.T) {
	for _, tc := range contentTypeTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}
