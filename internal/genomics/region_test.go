package genomics_test

import (
	"testing"

	. "github.com/david-xliu/htsget-refserver/internal/genomics"
)

func TestString(t *testing.T) {
	r := &Region{Name: "chr1", Start: "0", End: "100"}
	s := r.String()

	if s != "chr1:0-100" {
		t.Errorf("String was incorrect, got: %v, want: %v.", s, "chr1:0-100")
	}
}
