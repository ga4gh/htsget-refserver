package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

const EOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")
const EOF_LEN = int64(len(EOF))

const dataSource = "http://s3.amazonaws.com/czbiohub-tabula-muris/"
const testFile = "A1-B001176-3_56_F-1-1_R1.mus.Aligned.out.sorted.bam"

// Ticket holds the entire json ticket returned to the client
type ticket struct {
	HTSget container `json:"htsget"`
}

// Container holds the file format, urls of files for the client,
// and optionally an MD5 digest resulting from the concatenation of url data blocks
type container struct {
	Format string    `json:"format"`
	URLS   []urlJSON `json:"urls"`
	MD5    string    `json:"md5,omitempty"`
}

// URL holds the url, headers and class
type urlJSON struct {
	URL     string   `json:"url"`
	Headers *headers `json:"headers,omitempty"`
	Class   string   `json:"class,omitempty"`
}

// Headers contains any headers sent by the client such as
// the range of the query and authorization tokens
type headers struct {
	Range string `json:"range"`
}

var FIELDS map[string]int = map[string]int{
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

func getReads(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dataEndpoint := getDataURL()

	//send Head request to check that file exists and to get file size
	res, err := http.Head(dataSource + filePath(id))
	if err != nil {
		panic(err)
	}
	res.Body.Close()
	if res.Status == "404 Not Found" {
		fmt.Println(res.Status)
		// TODO return and send error message
		return
	}

	// *** Parse query params ***
	params := r.URL.Query()
	format, err := parseFormat(params)
	//reqClass, err := parseClass(params)
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	if err != nil {
		panic(err)
	}

	numBlocks := int(math.Floor((float64(numBytes) / (9 * math.Pow10(8)))))
	if numBlocks == 0 {
		numBlocks = 1
	}

	// build HTTP response
	u := make([]urlJSON, 0)
	var respClass string
	var h *headers
	if numBlocks == 1 {
		u = append(u, urlJSON{dataEndpoint.String(), h, respClass})
	} else {
		var start int64 = 0
		var blockSize int64 = int64(math.Floor((float64(numBytes) / float64(numBlocks))))
		var end int64 = start + blockSize
		for i := 1; i <= numBlocks; i++ {
			if i == 1 { // first of multiple blocks
				h = &headers{"bytes=" + "0" + "-" + strconv.FormatInt(hLen, 10)}
				start = hLen + 1
			} else {
				if end > numBytes {
					end = numBytes
				}
				h = &headers{"bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)}
			}
			u = append(u, urlJSON{dataEndpoint.String(), h, respClass})
		}
	}
	c := container{format, u, ""}
	t := ticket{HTSget: c}
	ticketJSON, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	// send back response
	w.Header().Set("Content-Type", "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8")
	w.Write(ticketJSON)
}

func filePath(id string) string {
	var path string
	if strings.HasPrefix(id, "10X") {
		path = "10x_bam_files/" + id
	} else {
		path = "facs_bam_files/" + id
	}
	return path
}

func getDataURL() *URL {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse("localhost:3000/data/")
	if err != nil {
		panic(err)
	}

	if os.Getenv("APP_ENV") == "production" {
		dataEndpoint.Path += id
	} else {
		dataEndpoint.Opaque += id
	}
}
