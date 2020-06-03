package htsgetutils

import (
	"testing"
)

func TestAddTrailingSlash(t *testing.T) {
	testCases := []*struct {
		url, expectedUrl string
	}{
		{"https://example.org", "https://example.org/"},
		{"https://htsget.ga4gh.org", "https://htsget.ga4gh.org/"},
		{"http://localhost:3000/", "http://localhost:3000/"},
		{"https://htsget.ga4gh.org/", "https://htsget.ga4gh.org/"},
	}

	for _, testCase := range testCases {
		if AddTrailingSlash(testCase.url) != testCase.expectedUrl {
			t.Errorf("url does not match expected when trailing slash added")
		}
	}
}
