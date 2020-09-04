package htsconstants

/* **************************************************
 * SERVICE INFO TYPE
 * ************************************************** */

var ServiceInfoTypeGroup = "org.ga4gh"

var ServiceInfoTypeArtifact = "htsget"

var ServiceInfoTypeVersion = "1.2.0"

/* **************************************************
 * COMMON SERVICE INFO DEFAULTS
 * ************************************************** */

var DfltServiceInfoOrganizationName = "Global Alliance for Genomics and Health"

var DfltServiceInfoOrganizationURL = "https://ga4gh.org"

var DfltServiceInfoContactURL = "mailto:jeremy.adams@ga4gh.org"

var DfltServiceInfoDocumentationURL = "https://ga4gh.org"

var DfltServiceInfoCreatedAt = StartupTime

var DfltServiceInfoUpdatedAt = StartupTime

var DfltServiceInfoEnvironment = "test"

var DfltServiceInfoVersion = "1.3.0"

/* **************************************************
 * READS-SPECIFIC DEFAULTS
 * ************************************************** */

var DfltServiceInfoReadsID = "htsgetref.reads"

var DfltServiceInfoReadsName = "GA4GH htsget reference server reads endpoint"

var DfltServiceInfoReadsDescription = "Stream alignment files (BAM/CRAM) according to GA4GH htsget protocol"

/* **************************************************
 * VARIANTS-SPECIFIC DEFAULTS
 * ************************************************** */

var DfltServiceInfoVariantsID = "htsgetref.variants"

var DfltServiceInfoVariantsName = "GA4GH htsget reference server variants endpoint"

var DfltServiceInfoVariantsDescription = "Stream variant files (VCF/BCF) according to GA4GH htsget protocol"

/* **************************************************
 * HTSGET EXTENSION CONSTANTS
 * ************************************************** */

var HtsgetExtensionDatatypeReads = "reads"

var HtsgetExtensionDatatypeVariants = "variants"
