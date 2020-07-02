// Package config allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults.go contains default runtime properties when not overriden
// by environment properties
package config

// getDefaults gets all default properties
//
// Returns
//	(map[string]string): map of default properties
func getDefaults() map[string]string {
	defaults := map[string]string{
		"port": "3000",
		"host": "http://localhost:3000",
	}
	return defaults
}

func getDefaultReadsSourcesRegistry() *DataSourceRegistry {
	sources := []map[string]string{
		{
			"pattern": "^10X",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{id[0:2]}",
		},
		{
			"pattern": "tabulamuris\\..*",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{id[0:2]}",
		},
	}

	registry := newDataSourceRegistry()
	for i := 0; i < len(sources); i++ {
		registry.addDataSource(newDataSource(sources[i]["pattern"], sources[i]["path"]))
	}
	return registry
}
