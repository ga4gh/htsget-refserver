package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

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

	// send Head request to check that file exists and to get file size
	/* res, err := http.Head("https://golang.org")*/
	//if err != nil {
	//panic(err)
	//}
	//defer res.Body.Close()
	//size := res.ContentLength
	//w.Write(res)

	// *** Parse query params ***
	params := r.URL.Query()
	format, err := parseFormat(params)
	class, err := parseClass(params)
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	_, err = parseFields(params)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("ContentLength:%v", contentlength)

	// The address of the endpoint on this server which serves the data
	dataEndpoint, err := url.Parse("localhost:3000/data/")
	if err != nil {
		panic(err)
	}

	if os.Getenv("APP_ENV") == "production" {
		dataEndpoint.Path += id
	} else {
		dataEndpoint.Opaque += id
	}
	// Add Query Parameters to the URL
	dataEndpoint.RawQuery = params.Encode() // Escape Query Parameters

	// build HTTP response
	h := &headers{"bytes=" + strconv.FormatUint(start, 10) + "-" + strconv.FormatUint(end, 10)}
	u := []urlJSON{{dataEndpoint.String(), h, class}}
	c := container{format, u, ""}
	t := ticket{HTSget: c}
	ticketJSON, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	// send back response
	w.Header().Set("Content-Type", "application/vnd.ga4gh.htsget.v1.0.0+json; charset=utf-8")
	w.Write(ticketJSON)
}

func parseFormat(params url.Values) (string, error) {
	if _, ok := params["format"]; ok {
		if validReadFormat(params["format"][0]) {
			return strings.ToUpper(params["format"][0]), nil
		} else {
			panic("UnsupportedFormat")
		}
	} else {
		return "BAM", nil
	}
}

func parseClass(params url.Values) (string, error) {
	if _, ok := params["class"]; ok {
		if validClass(params["class"][0]) {
			return strings.ToLower(params["class"][0]), nil
		} else {
			panic("InvalidInput")
		}
	}
	return "", nil
}

func parseRefName(params url.Values) (string, error) {
	if _, ok := params["referenceName"]; ok {
		return params["referenceName"][0], nil
	}
	return "", nil
}

func parseRange(params url.Values, refName string) (uint64, uint64, error) {
	if _, ok := params["start"]; ok {
		if _, ok := params["end"]; ok {
			if validRange(params["start"][0], params["end"][0], refName) {
				start, _ := strconv.ParseUint(params["start"][0], 10, 32)
				end, _ := strconv.ParseUint(params["end"][0], 10, 32)
				return start, end, nil
			} else {
				panic("InvalidRange")
			}
		}
	} else if _, ok := params["end"]; ok {
		panic("InvalidRange")
	}
	return 0, 0, nil
}

func parseFields(params url.Values) ([]string, error) {
	if _, ok := params["fields"]; ok {
		fields := strings.Split(params["fields"][0], ",")
		if !validFields(fields) {
			panic("InvalidInput")
		}
		return fields, nil
	}
	return []string{}, nil
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

func filePath(id string) string {
	var path string
	if strings.HasPrefix(id, "10X") {
		path = "10x_bam_files/" + id
	} else {
		path = "facs_bam_files/" + id
	}
	return path
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
	for _, field := range fields {
		if _, ok := FIELDS[field]; ok {
			return false
		}
	}
	return true
}
