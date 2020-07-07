// Package htsgetconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module environment.go contains operations for parsing runtime properties
// from various environment variables
package htsgetconfig

import (
	"os"
)

// getEnvironment retrieves valid properties set by environment variables
//
// Returns
//	(map[string]string): key-value map of properties set on environment
func getEnvironment() map[string]string {
	// environment variables the program scans for, and the configuration
	// property name they map to
	envConfigKeys := [2][2]string{
		{"HTSGET_PORT", "port"},
		{"HTSGET_HOST", "host"},
	}

	// scan for each environment variable, if it exists, add the value to the
	// map under the configuration property name
	environment := map[string]string{}
	for i := 0; i < len(envConfigKeys); i++ {
		envKey := envConfigKeys[i][0]
		configKey := envConfigKeys[i][1]
		value := os.Getenv(envKey)
		if len(value) != 0 {
			environment[configKey] = value
		}
	}
	return environment
}
