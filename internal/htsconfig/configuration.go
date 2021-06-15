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

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"

	"github.com/getlantern/deepcopy"
)

// Configuration contains properties loaded from the JSON config file
//
// Attributes
//	ReadsDataSourceRegistry (*DataSourceRegistry): data sources for reads endpoint
type Configuration struct {
	Container *configurationContainer `json:"htsgetConfig"`
}

type configurationContainer struct {
	ServerProps    *configurationServerProps `json:"props"`
	ReadsConfig    *configurationEndpoint    `json:"reads"`
	VariantsConfig *configurationEndpoint    `json:"variants"`
}

type configurationServerProps struct {
	Port                 string `json:"port"`
	Host                 string `json:"host"`
	DocsDir              string `json:"docsDir"`
	TempDir              string `json:"tempdir"`
	LogFile              string `json:"logFile"`
	CorsAllowedOrigins   string `json:"corsAllowedOrigins"`
	CorsAllowedMethods   string `json:"corsAllowedMethods"`
	CorsAllowedHeaders   string `json:"corsAllowedHeaders"`
	CorsAllowCredentials *bool  `json:"corsAllowCredentials"`
	CorsMaxAge           int    `json:"corsMaxAge"`
	AwsAssumeRole        *bool  `json:"awsAssumeRole"`
}

type configurationEndpoint struct {
	Enabled            *bool               `json:"enabled,true" default:"true"`
	DataSourceRegistry *DataSourceRegistry `json:"dataSourceRegistry"`
	ServiceInfo        *ServiceInfo        `json:"serviceInfo"`
}

var configurationSingleton *Configuration

var configurationSingletonLoaded = false

var configurationSingletonLoadedError error

func patchConfiguration(defR reflect.Value, patchR reflect.Value) {

	for i := 0; i < patchR.NumField(); i++ {
		defRType := defR.Field(i).Type().String()
		patchRType := patchR.Field(i).Type().String()
		var defRVal reflect.Value
		var patchRVal reflect.Value

		typesToPatch := []string{
			"string",
			"int",
			"*bool",
			"*htsconfig.DataSourceRegistry",
		}

		if !htsutils.IsItemInArray(defRType, typesToPatch) && !htsutils.IsItemInArray(patchRType, typesToPatch) {
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
			} else if defRType == "int" {
				defR.Field(i).Set(patchR.Field(i))
			} else if defRType == "*bool" {
				if !patchR.Field(i).IsNil() {
					defR.Field(i).Set(patchR.Field(i))
				}
			} else if defRType == "*htsconfig.DataSourceRegistry" {
				if !patchR.Field(i).IsNil() {
					defR.Field(i).Set(patchR.Field(i))
				}
			}
		}
	}
}

func LoadConfig() {
	newConfiguration := new(Configuration)
	deepcopy.Copy(newConfiguration, DefaultConfiguration)

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
	SetConfig(newConfiguration)
	configurationSingletonLoaded = true
}

func SetConfig(config *Configuration) {
	configurationSingleton = config
}

func GetConfig() *Configuration {
	if !configurationSingletonLoaded {
		LoadConfig()
	}
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

func SetHost(host string) {
	getServerProps().Host = host
}

// GetHost gets the current configuration 'host' setting, the host base url the
// service is running at
func GetHost() string {
	return htsutils.AddTrailingSlash(getServerProps().Host)
}

func GetDocsDir() string {
	return getServerProps().DocsDir
}

func GetTempDir() string {
	return htsutils.AddTrailingSlash(getServerProps().TempDir)
}

func GetTempFilePath(filename string) string {
	return filepath.Join(GetTempDir(), filename)
}

func CreateTempFile(filename string) (*os.File, error) {
	return os.Create(GetTempFilePath(filename))
}

func RemoveTempfile(file *os.File) error {
	return os.Remove(file.Name())
}

func GetLogFile() string {
	return getServerProps().LogFile
}

func GetCorsAllowedOrigins() string {
	return getServerProps().CorsAllowedOrigins
}

func GetCorsAllowedMethods() string {
	return getServerProps().CorsAllowedMethods
}

func GetCorsAllowedHeaders() string {
	return getServerProps().CorsAllowedHeaders
}

func GetCorsAllowCredentials() bool {
	return *getServerProps().CorsAllowCredentials
}

func GetCorsMaxAge() int {
	return getServerProps().CorsMaxAge
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

func IsAwsAssumeRole() bool {
	return *getServerProps().AwsAssumeRole
}
