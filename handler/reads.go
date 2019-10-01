package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

var dataSource = "http://s3.amazonaws.com/czbiohub-tabula-muris/"

// Ticket holds the entire json ticket returned to the client
type Ticket struct {
	HTSget Container `json:"htsget"`
}

// Container holds the file format, urls of files for the client,
// and optionally an MD5 digest resulting from the concatenation of url data blocks
type Container struct {
	Format string `json:"format"`
	URLS   []URL  `json:"urls"`
	MD5    string `json:"md5,omitempty"`
}

// URL holds the url, headers and class
type URL struct {
	URL     string   `json:"url"`
	Headers *Headers `json:"headers,omitempty"`
	Class   string   `json:"class,omitempty"`
}

// Headers contains any headers sent by the client such as
// the range of the query and authorization tokens
type Headers struct {
	Range string `json:"range"`
}

func getReads(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	// *** Parse query params ***
	params := req.URL.Query()

	// format param - optional
	var format string
	if _, ok := params["format"]; ok {
		if validReadFormat(params["format"][0]) {
			format = strings.ToUpper(params["format"][0])
		} else {
			panic("UnsupportedFormat")
		}
	} else {
		format = "BAM"
	}

	// class param
	var class string
	if _, ok := params["class"]; ok {
		if validClass(params["class"][0]) {
			class = strings.ToLower(params["class"][0])
		} else {
			panic("InvalidInput")
		}
	}

	// referenceName param
	var referenceName string
	if _, ok := params["referenceName"]; ok {
		referenceName = params["referenceName"][0]
	}

	// start/end params
	var start uint64
	var end uint64
	if _, ok := params["start"]; ok {
		if _, ok := params["end"]; ok {
			if validRange(params["start"][0], params["end"][0], referenceName) {
				start, _ = strconv.ParseUint(params["start"][0], 10, 32)
				end, _ = strconv.ParseUint(params["end"][0], 10, 32)
			} else {
				panic("InvalidRange")
			}
		}
	} else if _, ok := params["end"]; ok {
		panic("InvalidRange")
	}

	// fields params
	var fields []string
	if _, ok := params["fields"]; ok {
		fields = strings.Split(params["fields"][0], ",")
		if !validFields(fields) {
			panic("InvalidInput")
		}
	}

	var fileName string
	if strings.HasPrefix(id, "10X") {
		fileName = "10x_bam_files/" + id
	} else {
		fileName = "facs_bam_files/" + id
	}

	var md5 string
	var headers *Headers
	headers = &Headers{"bytes=" + strconv.FormatUint(start, 10) + "-" + strconv.FormatUint(end, 10)}
	urls := []URL{{dataSource + fileName, headers, class}}
	container := Container{format, urls, md5}
	ticket := Ticket{HTSget: container}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.ga4gh.htsget.v1.0.0+json; charset=utf-8")
	w.Write(ticketJSON)
}

func validReadFormat(s string) bool {
	switch strings.ToUpper(s) {
	case "BAM":
		return true
	case "CRAM":
		return true
	default:
		return false
	}
}

func validClass(s string) bool {
	switch strings.ToLower(s) {
	case "head":
		return true
	case "body":
		return true
	default:
		return false
	}
}

func validRange(startStr string, endStr string, refName string) bool {
	start, errStart := strconv.ParseUint(startStr, 10, 32)
	end, errEnd := strconv.ParseUint(endStr, 10, 32)

	if errStart != nil || errEnd != nil {
		return false
	}
	if start > end {
		return false
	}
	if refName == "" || refName == "*" {
		return false
	}

	return true
}

func validFields(fields []string) bool {
	fieldsMap := map[string]bool{
		"QNAME": true, // read names
		"FLAG":  true, // read bit flags
		"RNAME": true, // reference sequence name
		"POS":   true, // alignment position
		"MAPQ":  true, // mapping quality score
		"CIGAR": true, // CIGAR string
		"RNEXT": true, // reference sequence name of the next fragment template
		"PNEXT": true, // alignment position of the next fragment in the template
		"TLEN":  true, // inferred template size
		"SEQ":   true, // read bases
		"QUAL":  true, // base quality scores
	}

	for _, field := range fields {
		if !fieldsMap[field] {
			return false
		}
	}
	return true
}
