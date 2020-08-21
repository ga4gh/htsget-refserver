// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go amalgamates all runtime configuration sources
// (environment, JSON config file, cli, defaults) into a single object
package htsconfig

import (
	"errors"
	"strconv"
	"sync"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

// Configuration contains runtime properties loaded from env, config file, or default
//
// Attributes
// 	props (map[string]string): runtime properties dictionary
//	readsDataSourceRegistry (*DataSourceRegistry): data sources for /reads endpoint
type Configuration struct {
	props                      map[string]string
	readsDataSourceRegistry    *DataSourceRegistry
	variantsDataSourceRegistry *DataSourceRegistry
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
	config.readsDataSourceRegistry = getDefaultReadsSourcesRegistry()
	config.variantsDataSourceRegistry = getDefaultVariantsSourcesRegistry()
	if cliargs.configFile != "" {
		configFile := getConfigFile()
		configFileLoadError := getConfigFileLoadError()
		if configFileLoadError != nil {
			configLoadError = errors.New(configFileLoadError.Error())
			return
		}

		if configFile.ReadsDataSourceRegistry != nil {
			config.readsDataSourceRegistry = configFile.ReadsDataSourceRegistry
		}

		if configFile.VariantsDataSourceRegistry != nil {
			config.variantsDataSourceRegistry = configFile.VariantsDataSourceRegistry
		}
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
	return htsutils.AddTrailingSlash(getConfigProp("host"))
}

// GetReadsDataSourceRegistry gets the registered data sources for the 'reads' endpoint
//
// Returns
//	(*DataSourceRegistry): all 'reads' endpoint data sources
func GetReadsDataSourceRegistry() *DataSourceRegistry {
	return getConfig().readsDataSourceRegistry
}

func GetVariantsDataSourceRegistry() *DataSourceRegistry {
	return getConfig().variantsDataSourceRegistry
}

// getReadsPathForID gets a complete url or file path for a given ID
// given the request ID, this function looks up the 'reads' data source registry
// and finds the first data source matching the pattern. The id is then used to
// populate the path to the resource based on the data source's 'path' attribute
//
// Arguments
//	id (string): request ID
// Returns
//	(string): path to the object for the given id
//	(error): no match was found, or another error was encountered
func getReadsPathForID(id string) (string, error) {
	return GetReadsDataSourceRegistry().GetMatchingPath(id)
}

func getVariantsPathForID(id string) (string, error) {
	return GetVariantsDataSourceRegistry().GetMatchingPath(id)
}

func GetPathForID(endpoint htsconstants.ServerEndpoint, id string) (string, error) {
	functionsByEndpoint := [6]func(string) (string, error){
		getReadsPathForID,
		getReadsPathForID,
		nil,
		getVariantsPathForID,
		getVariantsPathForID,
		nil,
	}
	return functionsByEndpoint[endpoint](id)
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

	keys := [2]string{"readsDataSourceRegistry", "variantsDataSourceRegistry"}
	getters := [2]func() *DataSourceRegistry{GetReadsDataSourceRegistry, GetVariantsDataSourceRegistry}

	// validate 1. Reads, and 2. Variants data source registries coming from
	// config file if either is not set, use the default
	// if either is malformed, raise an error
	for i := 0; i < 2; i++ {
		registryFromConfig := getters[i]()

		configFileLoadError := getConfigFileLoadError()
		if configFileLoadError != nil {
			configFileLoadError = errors.New(configFileLoadError.Error())
			return
		}

		if registryFromConfig == nil {
			configLoadError = errors.New(keys[i] + " not configured, check json config file")
		}

		if registryFromConfig.Sources == nil {
			configLoadError = errors.New(keys[i] + " not configured, check json config file")
		}

		for j := 0; j < len(registryFromConfig.Sources); j++ {
			source := registryFromConfig.Sources[j]
			if source.Path == "" {
				msg := keys[i] + " incorrectly configured, missing \"path\" on source #" + strconv.Itoa(j)
				configLoadError = errors.New(msg)
				return

			}
			if source.Pattern == "" {
				msg := keys[i] + " incorrectly configured, missing \"pattern\" on source #" + strconv.Itoa(j)
				configLoadError = errors.New(msg)
				return
			}
		}
	}
}
