// Package htsticket produces the htsget JSON response ticket
//
// Module url_test tests url
package htsticket

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

var urlSetURLTC = []struct {
	url string
}{
	{"http://htsget.ga4gh.org/reads/data/object1"},
	{"http://localhost:3000/variants/data/1000genomes.00001"},
	{"http://localhost:4000/reads/data/gatktest.11111"},
}

var urlSetHeadersTC = []struct {
	blockid, numblocks, filepath string
}{
	{"1", "10", "./gatk/test1.bam"},
}

func TestUrlSetURL(t *testing.T) {
	for _, tc := range urlSetURLTC {
		url := NewURL()
		url.SetURL(tc.url)
		assert.Equal(t, tc.url, url.URL)
	}
}

func TestUrlSetHeaders(t *testing.T) {
	for _, tc := range urlSetHeadersTC {
		h := NewHeaders()
		h.SetBlockID(tc.blockid)
		h.SetNumBlocks(tc.numblocks)
		h.SetFilePathHeader(tc.filepath)
		url := NewURL()
		url.SetHeaders(h)
		assert.Equal(t, tc.blockid, url.Headers.BlockID)
		assert.Equal(t, tc.numblocks, url.Headers.NumBlocks)
		assert.Equal(t, tc.filepath, url.Headers.FilePath)
	}
}

func TestUrlSetClass(t *testing.T) {
	url := NewURL()
	functions := []func() *URL{
		url.SetClassHeader,
		url.SetClassBody,
	}
	exp := []string{
		htsconstants.ClassHeader,
		htsconstants.ClassBody,
	}
	for i := 0; i < len(functions); i++ {
		functions[i]()
		assert.Equal(t, exp[i], url.Class)
	}
}
