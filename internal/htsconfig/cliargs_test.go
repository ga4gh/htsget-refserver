// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module cliargs_test tests module cliargs
package htsconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var loadCliArgsTC = []struct {
	args       []string
	expCliArgs cliArgs
}{
	{
		[]string{"-config", "./data/config/config.json"},
		cliArgs{configFile: "./data/config/config.json"},
	},
}

func TestLoadCliArgs(t *testing.T) {

	for _, tc := range loadCliArgsTC {
		os.Args = append(os.Args, tc.args...)
		actualCliArgs := getCliArgs()
		assert.Equal(t, tc.expCliArgs.configFile, actualCliArgs.configFile)
	}
}
