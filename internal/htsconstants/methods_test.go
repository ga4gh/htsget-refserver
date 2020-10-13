// Package htsconstants contains program constants
//
// Module methods_test tests module methods
package htsconstants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// methodsTC test cases for Method enum
var methodsTC = []struct {
	e   HTTPMethod
	exp string
}{
	{GetMethod, "GET"},
	{PostMethod, "POST"},
}

// TestMethod tests Method enum
func TestMethod(t *testing.T) {
	for _, tc := range methodsTC {
		assert.Equal(t, tc.exp, tc.e.String())
	}
}
