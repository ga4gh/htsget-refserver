// Package config allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults.go contains default runtime properties when not overriden
// by environment properties
package config

// getDefaults gets all default properties
//
// Returns
//	(map[string]string): map of default properties
func getDefaults() map[string]string {
	defaults := map[string]string{
		"port": "3000",
		"host": "http://localhost:3000",
	}
	return defaults
}
