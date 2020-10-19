// Package htsticket produces the htsget JSON response ticket
//
// Module headers_test tests headers
package htsticket

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

// headersSetCurrentBlockTC test cases for SetCurrentBlock
var headersSetCurrentBlockTC = []struct {
	currentBlock string
}{
	{"1"}, {"2"}, {"3"},
}

// headersSetTotalBlocksTC test cases for SetTotalBlocks
var headersSetTotalBlocksTC = []struct {
	totalblocks string
}{
	{"10"}, {"20"}, {"100"},
}

// headersSetRangeHeaderTC test cases for SetRangeHeader
var headersSetRangeHeaderTC = []struct {
	start, end int64
	exp        string
}{
	{10000, 20000, "bytes=10000-20000"},
	{4567890, 9876543, "bytes=4567890-9876543"},
	{999, 1000, "bytes=999-1000"},
}

// headersSetFilePathHeaderTC test cases for SetFilePathHeader
var headersSetFilePathHeaderTC = []struct {
	filepath string
}{
	{"./data/gcp/gatk-test-data/wgs_bam/NA12878.bam"},
	{"./data/gcp/gatk-test-data/wgs_bam/NA12878_20k_b37.bam"},
}

// TestHeadersSetCurrentBlock tests SetCurrentBlock function
func TestHeadersSetCurrentBlock(t *testing.T) {
	for _, tc := range headersSetCurrentBlockTC {
		h := NewHeaders()
		h.SetCurrentBlock(tc.currentBlock)
		assert.Equal(t, tc.currentBlock, h.CurrentBlock)
	}
}

// TestHeadersSetNumBlocks tests SetTotalBlocks function
func TestHeadersSetTotalBlocks(t *testing.T) {
	for _, tc := range headersSetTotalBlocksTC {
		h := NewHeaders()
		h.SetTotalBlocks(tc.totalblocks)
		assert.Equal(t, tc.totalblocks, h.TotalBlocks)
	}
}

// TestHeadersSetRangeHeader tests SetRangeHeader function
func TestHeadersSetRangeHeader(t *testing.T) {
	for _, tc := range headersSetRangeHeaderTC {
		h := NewHeaders()
		h.SetRangeHeader(tc.start, tc.end)
		assert.Equal(t, tc.exp, h.Range)
	}
}

// TestHeadersSetClass tests SetClass function
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
		assert.Equal(t, exp[i], h.BlockClass)
	}
}

// TestHeadersSetFilePath tests SetFilePath function
func TestHeadersSetFilePath(t *testing.T) {
	for _, tc := range headersSetFilePathHeaderTC {
		h := NewHeaders()
		h.SetFilePathHeader(tc.filepath)
		assert.Equal(t, tc.filepath, h.FilePath)
	}
}
