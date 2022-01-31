// Package htsconstants contains program constants
//
// Module serviceinfo contains default properties for the service-info response
package htsconstants

/* **************************************************
 * SERVICE INFO TYPE
 * ************************************************** */

// ServiceInfoTypeGroup constant group value for service-info
var ServiceInfoTypeGroup = "org.ga4gh"

// ServiceInfoTypeArtifact constant artifact value for service-info
var ServiceInfoTypeArtifact = "htsget"

// ServiceInfoTypeVersion constant version value for service-info
var ServiceInfoTypeVersion = "1.2.0"

/* **************************************************
 * COMMON SERVICE INFO DEFAULTS
 * ************************************************** */

// DfltServiceInfoOrganizationName default organization name
var DfltServiceInfoOrganizationName = "Global Alliance for Genomics and Health"

// DfltServiceInfoOrganizationURL default organization url
var DfltServiceInfoOrganizationURL = "https://ga4gh.org"

// DfltServiceInfoContactURL default contact email / url
var DfltServiceInfoContactURL = "mailto:jeremy.adams@ga4gh.org"

// DfltServiceInfoDocumentationURL default documentation url
var DfltServiceInfoDocumentationURL = "https://ga4gh.org"

// DfltServiceInfoCreatedAt default created at time
var DfltServiceInfoCreatedAt = StartupTime

// DfltServiceInfoUpdatedAt default updated at time
var DfltServiceInfoUpdatedAt = StartupTime

// DfltServiceInfoEnvironment default environment
var DfltServiceInfoEnvironment = "test"

// DfltServiceInfoVersion default application version
var DfltServiceInfoVersion = "1.5.1"

/* **************************************************
 * READS-SPECIFIC DEFAULTS
 * ************************************************** */

// DfltServiceInfoReadsID default service-info id of reads API
var DfltServiceInfoReadsID = "htsgetref.reads"

// DfltServiceInfoReadsName default service-info name of reads API
var DfltServiceInfoReadsName = "GA4GH htsget reference server reads endpoint"

// DfltServiceInfoReadsDescription default service-info description of reads API
var DfltServiceInfoReadsDescription = "Stream alignment files (BAM/CRAM) according to GA4GH htsget protocol"

/* **************************************************
 * VARIANTS-SPECIFIC DEFAULTS
 * ************************************************** */

// DfltServiceInfoVariantsID default service-info id of variants API
var DfltServiceInfoVariantsID = "htsgetref.variants"

// DfltServiceInfoVariantsName default service-info name of variants API
var DfltServiceInfoVariantsName = "GA4GH htsget reference server variants endpoint"

// DfltServiceInfoVariantsDescription default service-info description of variants API
var DfltServiceInfoVariantsDescription = "Stream variant files (VCF/BCF) according to GA4GH htsget protocol"

/* **************************************************
 * HTSGET EXTENSION CONSTANTS
 * ************************************************** */

// HtsgetExtensionDatatypeReads datatype keyword for reads API
var HtsgetExtensionDatatypeReads = "reads"

// HtsgetExtensionDatatypeVariants datatype keyword for variants API
var HtsgetExtensionDatatypeVariants = "variants"
