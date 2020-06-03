package config

import "fmt"

func getEnvironment() map[string]string {
	environmentKeys := [1]string{
		"HTSGET_PORT",
	}

	environment := map[string]string{}

	for i := 0; i < len(environmentKeys); i++ {
		fmt.Println(environmentKeys[i])
	}

	return environment
}
