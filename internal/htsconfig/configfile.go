// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module configfile contains operations for setting properties from the
// JSON config file
package htsconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// configFile contains properties loaded from the JSON config file
//
// Attributes
//	ReadsDataSourceRegistry (*DataSourceRegistry): data sources for reads endpoint
type configFile struct {
	ReadsDataSourceRegistry    *DataSourceRegistry `json:"readsDataSourceRegistry"`
	VariantsDataSourceRegistry *DataSourceRegistry `json:"variantsDataSourceRegistry"`
}

// cfgFile (*configFile): singleton of config file settings
var cfgFile *configFile

// cfgFileLoad (sync.Once): indicates whether the singleton has been loaded or not
var cfgFileLoad sync.Once

// cfgFileLoadError (error): holds any error encountered during the setting of config file properties
var cfgFileLoadError error

// loadConfigFile instanties config file singleton with correct runtime properties
func loadConfigFile() {
	// get config file path from cli
	filePath := getCliArgs().configFile
	_, err := os.Stat(filePath)
	// check if the file doesn't exist, and if file is not valid JSON
	if os.IsNotExist(err) {
		cfgFileLoadError = errors.New("The specified config file doesn't exist: " + filePath)
		return
	}
	if err != nil {
		cfgFileLoadError = errors.New(err.Error())
		return
	}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		cfgFileLoadError = errors.New(err.Error())
		return
	}
	jsonContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		cfgFileLoadError = errors.New(err.Error())
		return
	}
	err = json.Unmarshal(jsonContent, &cfgFile)
	if err != nil {
		cfgFileLoadError = errors.New(err.Error())
	}
}

// getConfigFile get the the loaded configFile settings singleton
//
// Returns
//	(*configFile): loaded configFile singleton
func getConfigFile() *configFile {
	cfgFileLoad.Do(func() {
		loadConfigFile()
	})
	return cfgFile
}

// getConfigFileLoadError gets error object associated with config file loading
//
// Returns
//	(error): if set, an error was encountered during config file loading
func getConfigFileLoadError() error {
	return cfgFileLoadError
}
