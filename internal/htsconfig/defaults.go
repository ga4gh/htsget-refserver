// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module defaults contains default runtime properties when not overriden
// by environment properties
package htsconfig

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

var defaultServiceType = &ServiceType{
	Group:    htsconstants.ServiceInfoTypeGroup,
	Artifact: htsconstants.ServiceInfoTypeArtifact,
	Version:  htsconstants.ServiceInfoTypeVersion,
}

var defaultEnabledReads = true
var defaultFieldsParameterEffectiveReads = true
var defaultTagsParametersEffectiveReads = true

var defaultEnabledVariants = true
var defaultFieldsParameterEffectiveVariants = false
var defaultTagsParametersEffectiveVariants = false

var DefaultConfiguration = &Configuration{
	Container: &configurationContainer{
		ServerProps: &configurationServerProps{
			Port:                 htsconstants.DfltServerPropsPort,
			Host:                 htsconstants.DfltServerPropsHost,
			DocsDir:              htsconstants.DfltServerPropsDocsDir,
			TempDir:              htsconstants.DfltServerPropsTempDir,
			LogFile:              htsconstants.DfltServerPropsLogFile,
			LogFormat:            htsconstants.DfltServerPropsLogFormat,
			LogLevel:             htsconstants.DfltServerPropsLogLevel,
			CorsAllowedOrigins:   htsconstants.DfltCorsAllowedOrigins,
			CorsAllowedMethods:   htsconstants.DfltCorsAllowedMethods,
			CorsAllowedHeaders:   htsconstants.DfltCorsAllowedHeaders,
			CorsAllowCredentials: &htsconstants.DfltCorsAllowCredentials,
			CorsMaxAge:           htsconstants.DfltCorsMaxAge,
			ServerCert:           htsconstants.DfltServerPropsServerCert,
			ServerKey:            htsconstants.DfltServerPropsServerKey,
			AwsAssumeRole:        &htsconstants.DfltAwsAssumeRole,
		},
		ReadsConfig: &configurationEndpoint{
			Enabled: &defaultEnabledReads,
			DataSourceRegistry: &DataSourceRegistry{
				Sources: []*DataSource{
					&DataSource{
						Pattern: htsconstants.DfltReadsDataSourceTabulaMuris10XPattern,
						Path:    htsconstants.DfltReadsDataSourceTabulaMuris10XPath,
					},
					&DataSource{
						Pattern: htsconstants.DfltReadsDataSourceTabulaMurisFACSPattern,
						Path:    htsconstants.DfltReadsDataSourceTabulaMurisFACSPath,
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
				HtsgetExtension: &HtsgetExtension{
					Datatype:                 htsconstants.HtsgetExtensionDatatypeReads,
					Formats:                  htsconstants.APIEndpointReadsTicket.AllowedFormats(),
					FieldsParameterEffective: &defaultFieldsParameterEffectiveReads,
					TagsParametersEffective:  &defaultTagsParametersEffectiveReads,
				},
			},
		},
		VariantsConfig: &configurationEndpoint{
			Enabled: &defaultEnabledVariants,
			DataSourceRegistry: &DataSourceRegistry{
				Sources: []*DataSource{
					&DataSource{
						Pattern: htsconstants.DfltVariantsDataSource1000GPattern,
						Path:    htsconstants.DfltVariantsDataSource1000GPath,
					},
					&DataSource{
						Pattern: htsconstants.DfltVariantsDataSourceGIABTestPattern,
						Path:    htsconstants.DfltVariantsDataSourceGIABTestPath,
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
				HtsgetExtension: &HtsgetExtension{
					Datatype:                 htsconstants.HtsgetExtensionDatatypeVariants,
					Formats:                  htsconstants.APIEndpointVariantsTicket.AllowedFormats(),
					FieldsParameterEffective: &defaultFieldsParameterEffectiveVariants,
					TagsParametersEffective:  &defaultTagsParametersEffectiveVariants,
				},
			},
		},
	},
}
