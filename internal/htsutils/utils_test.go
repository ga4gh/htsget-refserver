// Package htsutils provides general, high-level, reusable functions
//
// Module utils_test tests utils
package htsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// utilsAddTrailingSlashTC test cases for AddTrailingSlash
var utilsAddTrailingSlashTC = []struct {
	url, exp string
}{
	{"https://example.org", "https://example.org/"},
	{"https://htsget.ga4gh.org", "https://htsget.ga4gh.org/"},
	{"http://localhost:3000/", "http://localhost:3000/"},
	{"https://htsget.ga4gh.org/", "https://htsget.ga4gh.org/"},
}

// utilsRemoveTrailingSlashTC test cases for RemoveTrailingSlash
var utilsRemoveTrailingSlashTC = []struct {
	url, exp string
}{
	{"https://example.org", "https://example.org"},
	{"https://htsget.ga4gh.org", "https://htsget.ga4gh.org"},
	{"http://localhost:3000/", "http://localhost:3000"},
	{"https://htsget.ga4gh.org/", "https://htsget.ga4gh.org"},
}

// utilsGetTagNameTC test cases for GetTagName
var utilsGetTagNameTC = []struct {
	tag, exp string
}{
	{"NH:i:1", "NH"},
	{"HI:i:1", "HI"},
	{"NM:i:0", "NM"},
	{"MD:Z:100", "MD"},
}

// utilsIsItemInArrayTC test cases for IsItemInArray
var utilsIsItemInArrayTC = []struct {
	item  string
	array []string
	exp   bool
}{
	{"NH", []string{"MD", "HI", "NM", "NH"}, true},
	{"NX", []string{"MD", "HI", "NM", "NH"}, false},
}

// utilsStringIsEmptyTC test cases for StringIsEmpty
var utilsStringIsEmptyTC = []struct {
	s   string
	exp bool
}{
	{"", true},
	{"empty", false},
	{"NH:i:1", false},
}

// utilsIsValidURLTC test cases for IsValidURL
var utilsIsValidURLTC = []struct {
	url string
	exp bool
}{
	{"http://localhost:3000/reads/object0", true},
	{"string", false},
	{"relative/path/to/object.bam", false},
}

// utilsParseRangeHeaderTC test cases for ParseRangeHeader
var utilsParseRangeHeaderTC = []struct {
	rangeHeader string
	expStart    int64
	expEnd      int64
	expErrorNil bool
}{
	{"bytes=10-50", 10, 50, true},
	{"malformedheader20to400", 0, 0, false},
	{"bytes=10.2-20.4", 0, 0, false},
}

// TestUtilsAddTrailingSlash tests AddTrailingSlash function
func TestUtilsAddTrailingSlash(t *testing.T) {
	for _, tc := range utilsAddTrailingSlashTC {
		assert.Equal(t, tc.exp, AddTrailingSlash(tc.url))
	}
}

// TestUtilsRemoveTrailingSlash tests RemoveTrailingSlash function
func TestUtilsRemoveTrailingSlash(t *testing.T) {
	for _, tc := range utilsRemoveTrailingSlashTC {
		assert.Equal(t, tc.exp, RemoveTrailingSlash(tc.url))
	}
}

// TestUtilsGetTagName tests GetTagName function
func TestUtilsGetTagName(t *testing.T) {
	for _, tc := range utilsGetTagNameTC {

		assert.Equal(t, tc.exp, GetTagName(tc.tag))
	}
}

// TestUtilsIsItemInArray tests IsItemInArray function
func TestUtilsIsItemInArray(t *testing.T) {
	for _, tc := range utilsIsItemInArrayTC {
		assert.Equal(t, tc.exp, IsItemInArray(tc.item, tc.array))
	}
}

// TestUtilsStringIsEmpty tests StringIsEmpty function
func TestUtilsStringIsEmpty(t *testing.T) {
	for _, tc := range utilsStringIsEmptyTC {
		assert.Equal(t, tc.exp, StringIsEmpty(tc.s))
	}
}

// TestUtilsIsValidUrl tests IsValidUrl function
func TestUtilsIsValidUrl(t *testing.T) {
	for _, tc := range utilsIsValidURLTC {
		assert.Equal(t, tc.exp, IsValidURL(tc.url))
	}
}

// TestUtilsParseRangeHeader tests ParseRangeHeader function
func TestUtilsParseRangeHeader(t *testing.T) {
	for _, tc := range utilsParseRangeHeaderTC {
		s, e, err := ParseRangeHeader(tc.rangeHeader)
		if tc.expErrorNil {
			assert.Empty(t, err)
			assert.Equal(t, tc.expStart, s)
			assert.Equal(t, tc.expEnd, e)
		} else {
			assert.NotEmpty(t, err)
		}
	}
}
