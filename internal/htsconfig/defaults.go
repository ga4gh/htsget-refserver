// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults contains default runtime properties when not overriden
// by environment properties
package htsconfig

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

var defaultEnabledReads = true
var defaultEnabledVariants = true
var defaultServiceType = &ServiceType{
	Group:    htsconstants.ServiceInfoTypeGroup,
	Artifact: htsconstants.ServiceInfoTypeArtifact,
	Version:  htsconstants.ServiceInfoTypeVersion,
}

var defaultConfiguration = &Configuration{
	Container: &configurationContainer{
		ServerProps: &configurationServerProps{
			Port: htsconstants.DfltServerPropsPort,
			Host: htsconstants.DfltServerPropsHost,
		},
		ReadsConfig: &configurationEndpoint{
			Enabled: &defaultEnabledReads,
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
			ServiceInfo: &ServiceInfo{
				ID:          htsconstants.DfltServiceInfoReadsID,
				Name:        htsconstants.DfltServiceInfoReadsName,
				Type:        defaultServiceType,
				Description: htsconstants.DfltServiceInfoReadsDescription,
				Organization: &Organization{
					Name: htsconstants.DfltServiceInfoOrganizationName,
					URL:  htsconstants.DfltServiceInfoOrganizationURL,
				},
				ContactURL:       htsconstants.DfltServiceInfoContactURL,
				DocumentationURL: htsconstants.DfltServiceInfoDocumentationURL,
				CreatedAt:        htsconstants.DfltServiceInfoCreatedAt,
				UpdatedAt:        htsconstants.DfltServiceInfoUpdatedAt,
				Environment:      htsconstants.DfltServiceInfoEnvironment,
				Version:          htsconstants.DfltServiceInfoVersion,
			},
		},
		VariantsConfig: &configurationEndpoint{
			Enabled: &defaultEnabledVariants,
			DataSourceRegistry: &DataSourceRegistry{
				Sources: []*DataSource{
					&DataSource{
						Pattern: "^1000genomes\\.(?P<accession>.*)$",
						Path:    "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz",
					},
				},
			},
			ServiceInfo: &ServiceInfo{
				ID:          htsconstants.DfltServiceInfoVariantsID,
				Name:        htsconstants.DfltServiceInfoVariantsName,
				Type:        defaultServiceType,
				Description: htsconstants.DfltServiceInfoVariantsDescription,
				Organization: &Organization{
					Name: htsconstants.DfltServiceInfoOrganizationName,
					URL:  htsconstants.DfltServiceInfoOrganizationURL,
				},
				ContactURL:       htsconstants.DfltServiceInfoContactURL,
				DocumentationURL: htsconstants.DfltServiceInfoDocumentationURL,
				CreatedAt:        htsconstants.DfltServiceInfoCreatedAt,
				UpdatedAt:        htsconstants.DfltServiceInfoUpdatedAt,
				Environment:      htsconstants.DfltServiceInfoEnvironment,
				Version:          htsconstants.DfltServiceInfoVersion,
			},
		},
	},
}
