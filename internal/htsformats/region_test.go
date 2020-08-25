// Package htsformats manipulates bioinformatic data encountered by htsget
//
// Module region_test tests region module
package htsformats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var regionStringTC = []struct {
	name, start, end, exp string
}{
	{"chr10", "-1", "-1", "chr10"},
	{"chr22", "100", "-1", "chr22:100"},
	{"chr5", "-1", "250000", "chr5:0-250000"},
	{"chr1", "0", "100", "chr1:0-100"},
}

var regionExportBcftoolsTC = []struct {
	name, start, end, exp string
}{
	{"chr10", "-1", "-1", "chr10"},
	{"chr22", "100", "-1", "chr22:100-"},
	{"chr5", "-1", "250000", "chr5:0-250000"},
	{"chr1", "0", "100", "chr1:0-100"},
}

func TestString(t *testing.T) {
	for _, tc := range regionStringTC {
		r := &Region{Name: tc.name, Start: tc.start, End: tc.end}
		assert.Equal(t, tc.exp, r.String())
	}
}

func TestExportSamtools(t *testing.T) {
	for _, tc := range regionStringTC {
		r := &Region{Name: tc.name, Start: tc.start, End: tc.end}
		assert.Equal(t, tc.exp, r.ExportSamtools())
	}
}

func TestExportBcftools(t *testing.T) {
	for _, tc := range regionExportBcftoolsTC {
		r := &Region{Name: tc.name, Start: tc.start, End: tc.end}
		assert.Equal(t, tc.exp, r.ExportBcftools())
	}
}
