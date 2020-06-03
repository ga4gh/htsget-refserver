package config

import (
	"sync"
)

type Configuration struct {
	props map[string]string
}

var config *Configuration
var once sync.Once

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
	config = newConfig
}

func getConfig() *Configuration {
	once.Do(func() {
		loadConfig()
	})
	return config
}

func GetConfigProp(key string) string {
	c := getConfig()
	return c.props[key]
}
