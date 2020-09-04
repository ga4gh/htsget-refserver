package htsrequest

import (
	"testing"
)

var parsePathParamTC = []struct {
	key, value, expString string
	expBool               bool
}{
	{"id", "tabulamuris.00001", "tabulamuris.00001", true},
}

func TestParsePathParam(t *testing.T) {
	// for _, tc := range parsePathParamTC {
	// }
}
