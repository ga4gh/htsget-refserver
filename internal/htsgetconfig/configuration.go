// Package htsgetconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go contains operations for setting runtime properties
// from the environment, config file, and defaults
package htsgetconfig

import (
	"errors"
	"strconv"
	"sync"

	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
)

// Configuration contains runtime properties loaded from env, config file, or default
//
// Attributes
// 	props (map[string]string): runtime properties dictionary
type Configuration struct {
	props                   map[string]string
	readsDataSourceRegistry *DataSourceRegistry
}

// config (Configuration): singleton of config to be used throughout the program
var config *Configuration

// configLoad (sync.Once): indicates whether the singleton config has been loaded or not
var configLoad sync.Once

var configLoadError error

// loadConfig instantiates config singleton with correct runtime properties
// loads properties into the map from defaults then environment (environment
// overrides defaults)
func loadConfig() {

	config = new(Configuration)
	config.props = make(map[string]string)
	defaults := getDefaults()
	environment := getEnvironment()
	// load properties dictionary
	// set default properties
	for k, v := range defaults {
		config.props[k] = v
	}
	// environment properties override defaults
	for k, v := range environment {
		config.props[k] = v
	}

	// load properties from cli args, then load reads data source registry
	cliargs := getCliArgs()
	readsDataSourcesRegistry := getDefaultReadsSourcesRegistry()
	if cliargs.configFile == "" {
		config.readsDataSourceRegistry = readsDataSourcesRegistry
	} else {
		configFile := getConfigFile()
		configFileLoadError := getConfigFileLoadError()
		if configFileLoadError != nil {
			configLoadError = errors.New(configFileLoadError.Error())
		}
		config.readsDataSourceRegistry = configFile.ReadsDataSourceRegistry
	}
}

// getConfig get the loaded config singleton
//
// Returns
// (*Configuration): loaded config singleton
func getConfig() *Configuration {
	configLoad.Do(func() {
		loadConfig()
	})
	return config
}

// getConfigProp get a single runtime property by its key
//
// Arguments
// 	key (string): property key
// Returns
//	(string): value for the specified property
func getConfigProp(key string) string {
	c := getConfig()
	return c.props[key]
}

func GetPort() string {
	return getConfigProp("port")
}

func GetHost() string {
	return htsgetutils.AddTrailingSlash(getConfigProp("host"))
}

func GetReadsDataSourceRegistry() *DataSourceRegistry {
	return getConfig().readsDataSourceRegistry
}

func GetReadsPathForID(id string) (string, error) {
	return GetReadsDataSourceRegistry().GetMatchingPath(id)
}

func GetConfigLoadError() error {
	return configLoadError
}

func LoadAndValidateConfig() {

	readsDataSourceRegistry := GetReadsDataSourceRegistry()
	if readsDataSourceRegistry == nil {
		configLoadError = errors.New("readsDataSourceRegistry not configured, check json config file")
		return
	}
	if readsDataSourceRegistry.Sources == nil {
		configLoadError = errors.New("readsDataSourceRegistry not configured, check json config file")
	}
	for i := 0; i < len(readsDataSourceRegistry.Sources); i++ {
		source := readsDataSourceRegistry.Sources[i]
		if source.Path == "" {
			msg := "readsDataSourceRegistry incorrectly configured, missing \"path\" on source #" + strconv.Itoa(i)
			configLoadError = errors.New(msg)
			return

		}
		if source.Pattern == "" {
			msg := "readsDataSourceRegistry incorrectly configured, missing \"pattern\" on source #" + strconv.Itoa(i)
			configLoadError = errors.New(msg)
			return
		}
	}
}
