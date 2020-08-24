// Package htsticket produces the htsget JSON response ticket
//
// Module headers_test tests headers
package htsticket

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

var headersSetBlockIDTC = []struct {
	blockID string
}{
	{"1"}, {"2"}, {"3"},
}

var headersSetNumBlocksTC = []struct {
	numblocks string
}{
	{"10"}, {"20"}, {"100"},
}

var headersSetRangeHeaderTC = []struct {
	start, end int64
	exp        string
}{
	{10000, 20000, "bytes=10000-20000"},
	{4567890, 9876543, "bytes=4567890-9876543"},
	{999, 1000, "bytes=999-1000"},
}

var headersSetFilePathHeaderTC = []struct {
	filepath string
}{
	{"./data/gcp/gatk-test-data/wgs_bam/NA12878.bam"},
	{"./data/gcp/gatk-test-data/wgs_bam/NA12878_20k_b37.bam"},
}

func TestHeadersSetBlockID(t *testing.T) {
	for _, tc := range headersSetBlockIDTC {
		h := NewHeaders()
		h.SetBlockID(tc.blockID)
		if h.BlockID != tc.blockID {
			t.Errorf("Expected: %s, Actual: %s", tc.blockID, h.BlockID)
		}
	}
}

func TestHeadersSetNumBlocks(t *testing.T) {
	for _, tc := range headersSetNumBlocksTC {
		h := NewHeaders()
		h.SetNumBlocks(tc.numblocks)
		if h.NumBlocks != tc.numblocks {
			t.Errorf("Expected: %s, Actual: %s", tc.numblocks, h.NumBlocks)
		}
	}
}

func TestHeadersSetRangeHeader(t *testing.T) {
	for _, tc := range headersSetRangeHeaderTC {
		h := NewHeaders()
		h.SetRangeHeader(tc.start, tc.end)
		if h.Range != tc.exp {
			t.Errorf("Expected: %s, Actual: %s", tc.exp, h.Range)
		}
	}
}

func TestHeadersSetClass(t *testing.T) {
	h := NewHeaders()
	functions := []func() *Headers{
		h.SetClassHeader,
		h.SetClassBody,
	}
	exp := []string{
		htsconstants.ClassHeader,
		htsconstants.ClassBody,
	}

	for i := 0; i < len(functions); i++ {
		functions[i]()
		if h.Class != exp[i] {
			t.Errorf("Expected: %s, Actual: %s", exp[i], h.Class)
		}
	}

}

func TestHeadersSetFilePath(t *testing.T) {
	for _, tc := range headersSetFilePathHeaderTC {
		h := NewHeaders()
		h.SetFilePathHeader(tc.filepath)
		if h.FilePath != tc.filepath {
			t.Errorf("Expected: %s, Actual: %s", tc.filepath, h.FilePath)
		}
	}
}
