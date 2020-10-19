// Package htsconstants contains program constants
//
// Module serverprops contains default properties for the serverprops configuration
package htsconstants

/* **************************************************
 * SERVER PROPS
 * ************************************************** */

// DfltServerPropsPort default port the server runs on
var DfltServerPropsPort = "3000"

// DfltServerPropsHost default hostname the server refers to when pointing to data endpoints
var DfltServerPropsHost = "http://localhost:3000/"

// DfltServerPropsDocsDir default static files directory
var DfltServerPropsDocsDir = ""

// DfltServerPropsTempDir default temporary file directory
var DfltServerPropsTempDir = "."

// DfltServerPropsLogFile default logfile to write logs
var DfltServerPropsLogFile = "htsget-refserver.log"

/* **************************************************
 * READS DATA SOURCE REGISTRY
 * ************************************************** */

// DfltReadsDataSourceTabulaMuris10XPattern regex pattern for tabula muris 10x ids
var DfltReadsDataSourceTabulaMuris10XPattern = "^tabulamuris\\.(?P<accession>10X.*)$"

// DfltReadsDataSourceTabulaMuris10XPath resolved path to tabula muris 10x files
var DfltReadsDataSourceTabulaMuris10XPath = "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam"

// DfltReadsDataSourceTabulaMurisFACSPattern regex pattern for tabula muris facs ids
var DfltReadsDataSourceTabulaMurisFACSPattern = "^tabulamuris\\.(?P<accession>.*)$"

// DfltReadsDataSourceTabulaMurisFACSPath resolved path to tabula muris FACS files
var DfltReadsDataSourceTabulaMurisFACSPath = "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam"

/* **************************************************
 * VARIANTS DATA SOURCE REGISTRY
 * ************************************************** */

// DfltVariantsDataSource1000GPattern regex pattern for 1000 genomes vcf ids
var DfltVariantsDataSource1000GPattern = "^1000genomes\\.(?P<accession>.*)$"

// DfltVariantsDataSource1000GPath resolved path to 1000 genomes vcf files
var DfltVariantsDataSource1000GPath = "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz"

// DfltVariantsDataSourceGIABTestPattern regex pattern for GIAB filtered test vcfs
var DfltVariantsDataSourceGIABTestPattern = "^(?P<accession>.*)_GIAB$"

// DfltVariantsDataSourceGIABTestPath resolved path to GIAB filtered test vcfs
var DfltVariantsDataSourceGIABTestPath = "./data/test/sources/giab/{accession}_GIAB.filtered.vcf.gz"
