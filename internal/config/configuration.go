// Package config allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go contains operations for setting runtime properties
// from the environment, config file, and defaults
package config

import (
	"sync"
)

// Configuration contains runtime properties loaded from env, config, or default
//
// Attributes
// 	props (map[string]string): runtime properties dictionary
type Configuration struct {
	props                   map[string]string
	readsDataSourceRegistry *DataSourceRegistry
}

// config (Configuration): singleton of config to be used throughout the program
var config *Configuration

// once (sync.Once): indicates whether the singleton config has been loaded or not
var once sync.Once

// loadConfig instantiates config singleton with correct runtime properties
// loads properties into the map from defaults then environment (environment
// overrides defaults)
func loadConfig() {

	var newConfig = new(Configuration)
	newConfig.props = make(map[string]string)
	defaults := getDefaults()
	environment := getEnvironment()
	for k, v := range defaults {
		newConfig.props[k] = v
	}
	for k, v := range environment {
		newConfig.props[k] = v
	}

	newConfig.readsDataSourceRegistry = getDefaultReadsSourcesRegistry()
	config = newConfig
}

// getConfig get the loaded config singleton
//
// Returns
// (*Configuration): loaded config singleton
func getConfig() *Configuration {
	once.Do(func() {
		loadConfig()
	})
	return config
}

// GetConfigProp get a single runtime property by its key
//
// Arguments
// 	key (string): property key
// Returns
//	(string): value for the specified property
func GetConfigProp(key string) string {
	c := getConfig()
	return c.props[key]
}

func GetReadsDataSourceRegistry() *DataSourceRegistry {
	return getConfig().readsDataSourceRegistry
}

func GetReadsPathForID(id string) (string, error) {
	return GetReadsDataSourceRegistry().GetMatchingPath(id)
}
