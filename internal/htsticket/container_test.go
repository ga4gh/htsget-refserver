// Package htsticket produces the htsget JSON response ticket
//
// Module container_test tests container
package htsticket

import "testing"

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
		if container.Format != expFormat[i] {
			t.Errorf("Expected: %s, Actual: %s", expFormat[i], container.Format)
		}
	}
}

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
			if container.URLS[i].URL != tc.urls[i] {
				t.Errorf("Expected: %s, Actual: %s", tc.urls[i], container.URLS[i].URL)
			}
		}
	}
}
