// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults contains default runtime properties when not overriden
// by environment properties
package htsconfig

var defaultConfiguration = &configuration{
	Container: &configurationContainer{
		ServerProps: &configurationServerProps{
			Port: "3000",
			Host: "http://localhost:3000/",
		},
		ReadsConfig: &configurationEndpoint{
			Enabled: true,
			DataSourceRegistry: &DataSourceRegistry{
				Sources: []*DataSource{
					&DataSource{
						Pattern: "^tabulamuris\\.(?P<accession>10X.*)$",
						Path:    "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam",
					},
					&DataSource{
						Pattern: "^tabulamuris\\.(?P<accession>.*)$",
						Path:    "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam",
					},
				},
			},
			ServiceInfo: &configurationServiceInfo{
				ID:   "htsgetref.reads",
				Name: "htsget reference server reads",
			},
		},
		VariantsConfig: &configurationEndpoint{
			Enabled: true,
			DataSourceRegistry: &DataSourceRegistry{
				Sources: []*DataSource{
					&DataSource{
						Pattern: "^1000genomes\\.(?P<accession>.*)$",
						Path:    "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz",
					},
				},
			},
			ServiceInfo: &configurationServiceInfo{
				ID:   "htsgetref.variants",
				Name: "htsget reference server variants",
			},
		},
	},
}
