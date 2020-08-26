// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go amalgamates all runtime configuration sources
// (environment, JSON config file, cli, defaults) into a single object
package htsconfig

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"

	"github.com/getlantern/deepcopy"
)

// configFile contains properties loaded from the JSON config file
//
// Attributes
//	ReadsDataSourceRegistry (*DataSourceRegistry): data sources for reads endpoint
type configuration struct {
	Container *configurationContainer `json:"htsgetconfig"`
}

type configurationContainer struct {
	ServerProps    *configurationServerProps `json:"props"`
	ReadsConfig    *configurationEndpoint    `json:"reads"`
	VariantsConfig *configurationEndpoint    `json:"variants"`
}

type configurationServerProps struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type configurationEndpoint struct {
	Enabled            bool                      `json:"enabled,true" default:"true"`
	DataSourceRegistry *DataSourceRegistry       `json:"dataSourceRegistry"`
	ServiceInfo        *configurationServiceInfo `json:"serviceInfo"`
}

type configurationServiceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var configurationSingleton *configuration

var configurationSingletonLoaded sync.Once

var configurationSingletonLoadedError error

func patchConfiguration(defR reflect.Value, patchR reflect.Value) {

	for i := 0; i < patchR.NumField(); i++ {
		defRType := defR.Field(i).Type().String()
		patchRType := patchR.Field(i).Type().String()
		var defRVal reflect.Value
		var patchRVal reflect.Value

		basicTypes := []string{"string", "bool"}

		if !htsutils.IsItemInArray(defRType, basicTypes) && !htsutils.IsItemInArray(patchRType, basicTypes) {
			defRVal = defR.Field(i).Elem()
			patchRVal = patchR.Field(i).Elem()

			defRValid := defRVal.IsValid()
			patchRValid := patchRVal.IsValid()
			if defRValid && patchRValid {
				patchConfiguration(defRVal, patchRVal)
			}
		} else {
			if defRType == "string" {
				patchString := patchR.Field(i).String()
				if patchString != "" {
					defR.Field(i).Set(patchR.Field(i))
				}
			} else if defRType == "bool" {
				test := patchR.Field(i).Bool()
				fmt.Println("test!")
				fmt.Println(test)
				fmt.Println("***")
				// patchBool := patchR.Field(i).Bool()
				// fmt.Println("patch bool")
				// fmt.Println(patchBool)
			}
		}
	}
}

func loadConfig() {
	newConfiguration := new(configuration)
	deepcopy.Copy(newConfiguration, defaultConfiguration)

	configFileLoadError := getConfigFileLoadError()
	if configFileLoadError != nil {
		configurationSingletonLoadedError = errors.New(configFileLoadError.Error())
	}

	configFileConfiguration := getConfigFile()
	patchConfiguration(
		reflect.ValueOf(newConfiguration).Elem(),
		reflect.ValueOf(configFileConfiguration).Elem(),
	)
	configurationSingleton = newConfiguration
}

func getConfig() *configuration {
	configurationSingletonLoaded.Do(func() {
		loadConfig()
	})
	return configurationSingleton
}

func getContainer() *configurationContainer {
	return getConfig().Container
}

func getServerProps() *configurationServerProps {
	return getContainer().ServerProps
}

// Port gets the current configuration 'port' setting, the port the server will run on
func Port() string {
	return getServerProps().Port
}

// Host gets the current configuration 'host' setting, the host base url the
// service is running at
func Host() string {
	return htsutils.AddTrailingSlash(getServerProps().Host)
}

func getEndpointConfig(ep htsconstants.APIEndpoint) *configurationEndpoint {
	reads := getContainer().ReadsConfig
	variants := getContainer().VariantsConfig
	configs := map[htsconstants.APIEndpoint]*configurationEndpoint{
		htsconstants.APIEndpointReadsTicket:         reads,
		htsconstants.APIEndpointReadsData:           reads,
		htsconstants.APIEndpointReadsServiceInfo:    reads,
		htsconstants.APIEndpointVariantsTicket:      variants,
		htsconstants.APIEndpointVariantsData:        variants,
		htsconstants.APIEndpointVariantsServiceInfo: variants,
	}
	return configs[ep]
}

func IsEndpointEnabled(ep htsconstants.APIEndpoint) bool {
	return getEndpointConfig(ep).Enabled
}

func getDataSourceRegistry(ep htsconstants.APIEndpoint) *DataSourceRegistry {
	return getEndpointConfig(ep).DataSourceRegistry
}

func GetObjectPath(ep htsconstants.APIEndpoint, id string) (string, error) {
	return getDataSourceRegistry(ep).GetMatchingPath(id)
}

/*
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
*/
