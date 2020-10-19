// Package htsticket produces the htsget JSON response ticket
//
// Module container_test tests container
package htsticket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// containerSetUrlsTC test cases for SetUrls
var containerSetUrlsTC = []struct {
	urls []string
}{
	{
		[]string{
			"http://htsget.ga4gh.org/reads/data/object1",
		},
	},
	{
		[]string{
			"http://htsget.ga4gh.org/reads/data/object1",
			"http://localhost:3000/variants/data/1000genomes.00001",
			"http://localhost:4000/reads/data/gatktest.11111",
		},
	},
}

// TestContainerSetFormat tests SetFormat function
func TestContainerSetFormat(t *testing.T) {
	container := NewContainer()
	functions := []func() *Container{
		container.SetFormatBam,
		container.SetFormatCram,
		container.SetFormatVcf,
		container.SetFormatBcf,
	}
	expFormat := []string{"BAM", "CRAM", "VCF", "BCF"}

	for i := 0; i < len(functions); i++ {
		functions[i]()
		assert.Equal(t, expFormat[i], container.Format)
	}
}

// TestContainerSetUrls tests SetUrls function
func TestContainerSetUrls(t *testing.T) {
	for _, tc := range containerSetUrlsTC {
		container := NewContainer()
		urls := []*URL{}

		for i := 0; i < len(tc.urls); i++ {
			url := NewURL()
			url.SetURL(tc.urls[i])
			urls = append(urls, url)
		}

		container.SetURLS(urls)

		for i := 0; i < len(tc.urls); i++ {
			assert.Equal(t, tc.urls[i], container.URLS[i].URL)
		}
	}
}
