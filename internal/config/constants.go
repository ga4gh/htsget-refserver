package config

import "encoding/hex"

var DATA_SOURCE_URL = "https://s3.amazonaws.com/czbiohub-tabula-muris/"

var N_BAM_FIELDS = 11

var BAM_FIELDS map[string]int = map[string]int{
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
var BAM_EXCLUDED_VALUES []string = []string{
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

var BAM_EOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")
var BAM_EOF_LEN = len(BAM_EOF)
var BAM_HEADER_EOF_LEN = 12
