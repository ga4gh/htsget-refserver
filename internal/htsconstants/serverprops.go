package htsconstants

/* **************************************************
 * SERVER PROPS
 * ************************************************** */

var DfltServerPropsPort = "3000"

var DfltServerPropsHost = "http://localhost:3000/"

var DfltServerPropsTempdir = "."

var DfltServerPropsLogfile = "htsget-refserver.log"

/* **************************************************
 * READS DATA SOURCE REGISTRY
 * ************************************************** */

var DfltReadsDataSourceTabulaMuris10XPattern = "^tabulamuris\\.(?P<accession>10X.*)$"

var DfltReadsDataSourceTabulaMuris10XPath = "https://s3.amazonaws.com/czbiohub-tabula-muris/10x_bam_files/{accession}_possorted_genome.bam"

var DfltReadsDataSourceTabulaMurisFACSPattern = "^tabulamuris\\.(?P<accession>.*)$"

var DfltReadsDataSourceTabulaMurisFACSPath = "https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/{accession}.mus.Aligned.out.sorted.bam"

/* **************************************************
 * VARIANTS DATA SOURCE REGISTRY
 * ************************************************** */

var DfltVariantsDataSource1000GPattern = "^1000genomes\\.(?P<accession>.*)$"

var DfltVariantsDataSource1000GPath = "https://ftp-trace.ncbi.nih.gov/1000genomes/ftp/phase1/analysis_results/integrated_call_sets/{accession}.vcf.gz"
