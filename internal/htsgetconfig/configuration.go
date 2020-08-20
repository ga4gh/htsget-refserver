// Package htsgetconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go amalgamates all runtime configuration sources
// (environment, JSON config file, cli, defaults) into a single object
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
//	readsDataSourceRegistry (*DataSourceRegistry): data sources for /reads endpoint
type Configuration struct {
	props                   map[string]string
	readsDataSourceRegistry *DataSourceRegistry
}

// config (*Configuration): singleton of config to be used throughout the program
var config *Configuration

// configLoad (sync.Once): indicates whether the singleton config has been loaded or not
var configLoad sync.Once

// configLoadError (error): holds any error encountered during setting of overall config
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
			return
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

// GetPort gets the current configuration 'port' setting
//
// Returns
//	(string): current port setting - the port the server will run on
func GetPort() string {
	return getConfigProp("port")
}

// GetHost gets the current configuration 'host' setting
//
// Returns
//	(string): host setting - the host base url the service is running at
func GetHost() string {
	return htsgetutils.AddTrailingSlash(getConfigProp("host"))
}

// GetReadsDataSourceRegistry gets the registered data sources for the 'reads' endpoint
//
// Returns
//	(*DataSourceRegistry): all 'reads' endpoint data sources
func GetReadsDataSourceRegistry() *DataSourceRegistry {
	return getConfig().readsDataSourceRegistry
}

// GetReadsPathForID gets a complete url or file path for a given ID
// given the request ID, this function looks up the 'reads' data source registry
// and finds the first data source matching the pattern. The id is then used to
// populate the path to the resource based on the data source's 'path' attribute
//
// Arguments
//	id (string): request ID
// Returns
//	(string): path to the object for the given id
//	(error): no match was found, or another error was encountered
func GetReadsPathForID(id string) (string, error) {
	return GetReadsDataSourceRegistry().GetMatchingPath(id)
}

// GetConfigLoadError gets the error associated with loading the configuration
//
// Returns
//	(error): if not nil, an error was encountered during configuration loading
func GetConfigLoadError() error {
	return configLoadError
}

// LoadAndValidateConfig performs custom validation on the configuration,
// ensuring properties are set correctly. sets the configuration error if
// any errors encountered
func LoadAndValidateConfig() {
	// validate reads data source registry is set and populated with valid
	// data sources
	readsDataSourceRegistry := GetReadsDataSourceRegistry()
	configFileLoadError := getConfigFileLoadError()
	if configFileLoadError != nil {
		configFileLoadError = errors.New(configFileLoadError.Error())
		return
	}
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
