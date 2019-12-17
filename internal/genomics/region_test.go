package genomics_test

import (
	"testing"

	. "github.com/david-xliu/htsget-refserver/internal/genomics"
)

func TestString(t *testing.T) {
	r := &Region{Name: "chr1", Start: "0", End: "100"}
	s := r.String()
	if s != "chr1:0-100" {
		t.Errorf("Got: %v, wanted: chr1:0-100", s)
	}

	r = &Region{Name: "chr1", Start: "-1"}
	s = r.String()
	if s != "chr1" {
		t.Errorf("Got: %v, wanted: chr1", s)
	}

	r = &Region{Name: "chr1", Start: "100", End: "-1"}
	s = r.String()
	if s != "chr1:100" {
		t.Errorf("Got: %v, wanted: chr1:100", s)
	}
}
