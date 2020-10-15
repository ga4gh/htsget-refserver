// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module region_test tests region module
package htsrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// regionTC test case struct for multiple region testing functions
type regionTC struct {
	nilStart, nilEnd bool
	referenceName    string
	start, end       int
	exp              string
}

// regionStringTC test cases for String
var regionStringTC = []regionTC{
	{false, false, "chr10", -1, -1, "chr10"},
	{false, false, "chr22", 100, -1, "chr22:100"},
	{false, false, "chr5", -1, 250000, "chr5:0-250000"},
	{false, false, "chr1", 0, 100, "chr1:0-100"},
	{true, false, "chr21", 0, 100000, "chr21:0-100000"},
	{false, true, "chr21", 100000, 0, "chr21:100000"},
}

// regionExportBcftoolsTC test cases for ExportBcftools
var regionExportBcftoolsTC = []regionTC{
	{false, false, "chr10", -1, -1, "chr10"},
	{false, false, "chr22", 100, -1, "chr22:100-"},
	{false, false, "chr5", -1, 250000, "chr5:0-250000"},
	{false, false, "chr1", 0, 100, "chr1:0-100"},
}

// regionReferenceNameTC test cases for GetReferenceName
var regionReferenceNameTC = []struct {
	referenceName string
	expRequested  bool
}{
	{"", false},
	{"chr1", true},
	{"chr22", true},
}

// instantiateRegion convenience method to construct a region, sometime with
// nil start/end if requested
func instantiateRegion(tc *regionTC) *Region {
	var r *Region

	if tc.nilStart || tc.nilEnd {
		var start *int
		var end *int
		if tc.nilStart {
			start = nil
		} else {
			start = &tc.start
		}
		if tc.nilEnd {
			end = nil
		} else {
			end = &tc.end
		}
		r = &Region{ReferenceName: tc.referenceName, Start: start, End: end}
	} else {
		r = NewRegion()
		r.SetReferenceName(tc.referenceName)
		r.SetStart(tc.start)
		r.SetEnd(tc.end)
	}
	return r
}

// TestRegionString tests String function
func TestRegionString(t *testing.T) {
	for _, tc := range regionStringTC {
		r := instantiateRegion(&tc)
		assert.Equal(t, tc.exp, r.String())
	}
}

// TestRegionExportSamtools tests ExportSamtools function
func TestRegionExportSamtools(t *testing.T) {
	for _, tc := range regionStringTC {
		r := instantiateRegion(&tc)
		assert.Equal(t, tc.exp, r.ExportSamtools())
	}
}

// TestRegionExportBcftools tests ExportBcftools function
func TestRegionExportBcftools(t *testing.T) {
	for _, tc := range regionExportBcftoolsTC {
		r := instantiateRegion(&tc)
		assert.Equal(t, tc.exp, r.ExportBcftools())
	}
}

// TestRegionGetReferenceName tests GetReferenceName function
func TestRegionGetReferenceName(t *testing.T) {
	for _, tc := range regionReferenceNameTC {
		r := &Region{ReferenceName: tc.referenceName, Start: nil, End: nil}
		assert.Equal(t, tc.referenceName, r.GetReferenceName())
		assert.Equal(t, tc.expRequested, r.ReferenceNameRequested())
	}
}
