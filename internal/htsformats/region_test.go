// Package htsformats manipulates bioinformatic data encountered by htsget
//
// Module region_test tests region module
package htsformats

import (
	"testing"
)

func TestString(t *testing.T) {
	r := &Region{Name: "chr1", Start: "0", End: "100"}
	s := r.String()
	if s != "chr1:0-100" {
		t.Errorf("Got: %v, expected: chr1:0-100", s)
	}

	r = &Region{Name: "chr1", Start: "-1"}
	s = r.String()
	if s != "chr1" {
		t.Errorf("Got: %v, expected: chr1", s)
	}

	r = &Region{Name: "chr1", Start: "100", End: "-1"}
	s = r.String()
	if s != "chr1:100" {
		t.Errorf("Got: %v, expected: chr1:100", s)
	}
}
