package config

import "encoding/hex"

var DATA_SOURCE_URL = "https://s3.amazonaws.com/czbiohub-tabula-muris/"

var BAM_FIELDS map[string]int = map[string]int{
	"QNAME": 1,  // read names
	"FLAG":  2,  // read bit flags
	"RNAME": 3,  // reference sequence name
	"POS":   4,  // alignment position
	"MAPQ":  5,  // mapping quality score
	"CIGAR": 6,  // CIGAR string
	"RNEXT": 7,  // reference sequence name of the next fragment template
	"PNEXT": 8,  // alignment position of the next fragment in the template
	"TLEN":  9,  // inferred template size
	"SEQ":   10, // read bases
	"QUAL":  11, // base quality scores
}
var BAM_EOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")
var BAM_EOF_LEN = len(BAM_EOF)
var BAM_HEADER_EOF_LEN = 12
