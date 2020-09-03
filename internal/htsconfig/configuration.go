// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configuration.go amalgamates all runtime configuration sources
// (environment, JSON config file, cli, defaults) into a single object
package htsconfig

import (
	"errors"
	"os"
	"path/filepath"
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
type Configuration struct {
	Container *configurationContainer `json:"htsgetconfig"`
}

type configurationContainer struct {
	ServerProps    *configurationServerProps `json:"props"`
	ReadsConfig    *configurationEndpoint    `json:"reads"`
	VariantsConfig *configurationEndpoint    `json:"variants"`
}

type configurationServerProps struct {
	Port    string `json:"port"`
	Host    string `json:"host"`
	Tempdir string `json:"tempdir"`
	Logfile string `json:"logfile"`
}

type configurationEndpoint struct {
	Enabled            *bool               `json:"enabled,true" default:"true"`
	DataSourceRegistry *DataSourceRegistry `json:"dataSourceRegistry"`
	ServiceInfo        *ServiceInfo        `json:"serviceInfo"`
}

var configurationSingleton *Configuration

var configurationSingletonLoaded sync.Once

var configurationSingletonLoadedError error

func patchConfiguration(defR reflect.Value, patchR reflect.Value) {

	for i := 0; i < patchR.NumField(); i++ {
		defRType := defR.Field(i).Type().String()
		patchRType := patchR.Field(i).Type().String()
		var defRVal reflect.Value
		var patchRVal reflect.Value

		basicTypes := []string{"string", "*bool"}

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
			} else if defRType == "*bool" {
				if !patchR.Field(i).IsNil() {
					defR.Field(i).Set(patchR.Field(i))
				}
			}
		}
	}
}

func loadConfig() {
	newConfiguration := new(Configuration)
	deepcopy.Copy(newConfiguration, defaultConfiguration)

	configFileLoadError := getConfigFileLoadError()
	if configFileLoadError != nil {
		configurationSingletonLoadedError = errors.New(configFileLoadError.Error())
	}

	configFileConfiguration := getConfigFile()
	if configFileConfiguration != nil {
		patchConfiguration(
			reflect.ValueOf(newConfiguration).Elem(),
			reflect.ValueOf(configFileConfiguration).Elem(),
		)
	}
	configurationSingleton = newConfiguration
}

func GetConfig() *Configuration {
	configurationSingletonLoaded.Do(func() {
		loadConfig()
	})
	return configurationSingleton
}

func getContainer() *configurationContainer {
	return GetConfig().Container
}

func getServerProps() *configurationServerProps {
	return getContainer().ServerProps
}

// GetPort gets the current configuration 'port' setting, the port the server will run on
func GetPort() string {
	return getServerProps().Port
}

// GetHost gets the current configuration 'host' setting, the host base url the
// service is running at
func GetHost() string {
	return htsutils.AddTrailingSlash(getServerProps().Host)
}

func GetTempdir() string {
	return htsutils.AddTrailingSlash(getServerProps().Tempdir)
}

func GetTempfilePath(filename string) string {
	return filepath.Join(GetTempdir(), filename)
}

func CreateTempfile(filename string) (*os.File, error) {
	return os.Create(GetTempfilePath(filename))
}

func RemoveTempfile(file *os.File) error {
	return os.Remove(file.Name())
}

func GetLogfile() string {
	return getServerProps().Logfile
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
	return *getEndpointConfig(ep).Enabled
}

func GetDataSourceRegistry(ep htsconstants.APIEndpoint) *DataSourceRegistry {
	return getEndpointConfig(ep).DataSourceRegistry
}

func GetObjectPath(ep htsconstants.APIEndpoint, id string) (string, error) {
	return GetDataSourceRegistry(ep).GetMatchingPath(id)
}

func GetServiceInfo(ep htsconstants.APIEndpoint) *ServiceInfo {
	return getEndpointConfig(ep).ServiceInfo
}

// GetConfigLoadError gets the error associated with loading the configuration
func GetConfigLoadError() error {
	return configurationSingletonLoadedError
}
