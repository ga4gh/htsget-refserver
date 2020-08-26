// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module environment_test tests module environment
package htsconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var environmentTC = []struct {
	environmentKey, dictionaryKey, value string
}{
	{"HTSGET_PORT", "port", "4000"},
	{"HTSGET_PORT", "port", "8989"},
	{"HTSGET_HOST", "host", "https://htsget.ga4gh.org/"},
	{"HTSGET_HOST", "host", "https://htsget-service-ga4gh.com/v1/"},
}

func TestEnvironment(t *testing.T) {
	for _, tc := range environmentTC {
		os.Setenv(tc.environmentKey, tc.value)
		env := getEnvironment()
		assert.Equal(t, tc.value, env[tc.dictionaryKey])
	}
}
