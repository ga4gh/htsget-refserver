// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults contains default runtime properties when not overriden
// by environment properties
package htsconfig

// getDefaults gets default properties as a map of strings
func getDefaults() map[string]string {
	defaults := map[string]string{
		"port": "3000",
		"host": "http://localhost:3000",
	}
	return defaults
}

// getDefaultReadsSourcesRegistry gets the default source registry for 'reads' endpoint
func getDefaultReadsSourcesRegistry() *DataSourceRegistry {
	sources := []map[string]string{
		{
			"pattern": "^tabulamuris\\.(?P<accession>10X.*)$",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam",
		},
		{
			"pattern": "^tabulamuris\\.(?P<accession>.*)$",
			"path":    "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam",
		},
	}

	registry := newDataSourceRegistry()
	for i := 0; i < len(sources); i++ {
		registry.addDataSource(newDataSource(sources[i]["pattern"], sources[i]["path"]))
	}
	return registry
}

// getDefaultVariantsSourcesRegistry gets the default source registry for 'variants' endpoint
func getDefaultVariantsSourcesRegistry() *DataSourceRegistry {
	sources := []map[string]string{
		{
			"pattern": "^1000genomes\\.(?P<accession>.*)$",
			"path":    "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz",
		},
	}
	registry := newDataSourceRegistry()
	for i := 0; i < len(sources); i++ {
		registry.addDataSource(newDataSource(sources[i]["pattern"], sources[i]["path"]))
	}
	return registry
}
