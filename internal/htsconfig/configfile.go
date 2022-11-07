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
)

var configFileSingleton *Configuration

var configFileSingletonLoaded = false

var configFileSingletonLoadedError error

// LoadConfigFile instantiates config file singleton with correct runtime properties
func LoadConfigFile() {
	// get config file path from cli
	filePath := getCliArgs().configFile
	_, err := os.Stat(filePath)
	// check if the file doesn't exist, and if file is not valid JSON
	if os.IsNotExist(err) {
		configFileSingletonLoadedError = errors.New("The specified config file doesn't exist: " + filePath)
		return
	}
	if err != nil {
		configFileSingletonLoadedError = errors.New(err.Error())
		return
	}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		configFileSingletonLoadedError = errors.New(err.Error())
		return
	}
	jsonContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		configFileSingletonLoadedError = errors.New(err.Error())
		return
	}

	err = json.Unmarshal(jsonContent, &configFileSingleton)
	if err != nil {
		configFileSingletonLoadedError = errors.New(err.Error())
	}
	configFileSingletonLoaded = true
}

func SetConfigFile(configFile *Configuration) {
	configFileSingleton = configFile
}

// getConfigFile get the the loaded configFile settings singleton
func getConfigFile() *Configuration {
	if !configFileSingletonLoaded {
		LoadConfigFile()
	}
	return configFileSingleton
}

// getConfigFileLoadError gets error object associated with config file loading
func getConfigFileLoadError() error {
	return configFileSingletonLoadedError
}
