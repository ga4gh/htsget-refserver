package config

func getDefaults() map[string]string {
	defaults := map[string]string{
		"port": "3000",
		"host": "http://localhost:3000",
	}
	return defaults
}
