package config

import (
	"os"
)

func getEnvironment() map[string]string {
	envConfigKeys := [2][2]string{
		{"HTSGET_PORT", "port"},
		{"HTSGET_HOST", "host"},
	}

	environment := map[string]string{}

	for i := 0; i < len(envConfigKeys); i++ {
		envKey := envConfigKeys[i][0]
		configKey := envConfigKeys[i][1]
		value := os.Getenv(envKey)
		if len(value) != 0 {
			environment[configKey] = value
		}
	}
	return environment
}
