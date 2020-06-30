// Package config allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module constants.go contains program constants
package config

import "encoding/hex"

// DataSourceURL (string): base url to AWS S3 bucket
var DataSourceURL = "https://s3.amazonaws.com/czbiohub-tabula-muris/"

// BamFieldsN (int): canonical number of fields in SAM/BAM (excluding tags)
var BamFieldsN = 11

// BamFields (map[string]int): ordered map of canonical column name to position
var BamFields map[string]int = map[string]int{
	"QNAME": 0,  // read names
	"FLAG":  1,  // read bit flags
	"RNAME": 2,  // reference sequence name
	"POS":   3,  // alignment position
	"MAPQ":  4,  // mapping quality score
	"CIGAR": 5,  // CIGAR string
	"RNEXT": 6,  // reference sequence name of the next fragment template
	"PNEXT": 7,  // alignment position of the next fragment in the template
	"TLEN":  8,  // inferred template size
	"SEQ":   9,  // read bases
	"QUAL":  10, // base quality scores
}

// BamExcludedValues ([]string): correct values when column is removed by column
var BamExcludedValues []string = []string{
	"*",   // QNAME
	"0",   // FLAG
	"*",   // RNAME
	"0",   // POS
	"255", // MAPQ
	"*",   // CIGAR
	"*",   // RNEXT
	"0",   // PNEXT
	"0",   // TLEN
	"*",   // SEQ
	"*",   // QUAL
}

// BamEOF ([]byte): BAM end of file byte sequence
var BamEOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")

// BamEOFLen (int): length (number of bytes) of BAM end of file byte sequence
var BamEOFLen = len(BamEOF)

// BamHeaderEOFLen (int): length (number of bytes) of BAM header end marker
var BamHeaderEOFLen = 12
