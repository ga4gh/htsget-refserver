package htsgetconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type configFile struct {
	ReadsDataSourceRegistry *DataSourceRegistry `json:"readsDataSourceRegistry"`
}

var cfgFile *configFile

var cfgFileLoad sync.Once

var cfgFileLoadError error

func loadConfigFile() {
	filePath := getCliArgs().configFile
	_, err := os.Stat(filePath)
	if err != nil {
		cfgFileLoadError = errors.New(err.Error())
		return
	}
	if os.IsNotExist(err) {
		cfgFileLoadError = errors.New("The specified config file doesn't exist: " + filePath)
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

func getConfigFile() *configFile {
	cfgFileLoad.Do(func() {
		loadConfigFile()
	})
	return cfgFile
}

func getConfigFileLoadError() error {
	return cfgFileLoadError
}
