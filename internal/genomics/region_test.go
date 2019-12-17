package genomics_test

import (
	"testing"

	. "github.com/david-xliu/htsget-refserver/internal/genomics"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	r := &Region{Name: "chr1", Start: "0", End: "100"}
	s := r.String()
	assert.Equal(t, "chr1:0-100", s, "They should be equal")

	r = &Region{Name: "chr1", Start: "-1"}
	s = r.String()
	assert.Equal(t, "chr1", s, "They should be equal")

	r = &Region{Name: "chr1", Start: "100", End: "-1"}
	s = r.String()
	assert.Equal(t, "chr1:100", s, "They should be equal")
}
