package server

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/genomics"
	"github.com/go-chi/chi"
)

var EOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")
var EOF_LEN = len(EOF)
var HEADER_EOF_LEN = 12

var dataSource = "http://s3.amazonaws.com/czbiohub-tabula-muris/"

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

// Headers contains any headers needed by the server from the client
type headers struct {
	BlockID   string `json:"block-id,omitempty"`   // id of current block
	NumBlocks string `json:"num-blocks,omitempty"` // total number of blocks
	Range     string `json:"range,omitempty"`
	Class     string `json:"class,omitempty"`
}

// htsgetError holds errors defined in the htsget protocol
type htsgetError struct {
	Code   int
	Htsget errorContainer `json:"htsget"`
}

type errorContainer struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (err *htsgetError) Error() string {
	return fmt.Sprint(err.Htsget.Error + ": " + err.Htsget.Message)
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
	var host string
	if os.Getenv("APP_ENV") == "production" {
		host = "https://htsget.ga4gh.org/"
	} else {
		host = "http://localhost:3000/"
	}

	//send Head request to check that file exists and to get file size
	res, err := http.Head(dataSource + filePath(id))
	if err != nil {
		writeError(w, err)
		return
	}
	res.Body.Close()
	if res.Status == "404 Not Found" {
		htsErr := &htsgetError{
			Code: http.StatusNotFound,
			Htsget: errorContainer{
				"NotFound",
				"The resource requested was not found",
			},
		}
		writeError(w, htsErr)
		return
	}

	// *** Parse query params ***
	params := r.URL.Query()
	format, err := parseFormat(params)
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"UnsupportedFormat",
				"The requested file format is not supported by the server",
			},
		}
		writeError(w, htsErr)
		return
	}

	queryClass, err := parseQueryClass(params)
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"InvalidInput",
				"The request parameters do not adhere to the specification",
			},
		}
		writeError(w, htsErr)
		return
	}

	refName := parseRefName(params)
	start, end, err := parseRange(params, refName)
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"InvalidRange",
				"The request range cannot be satisfied",
			},
		}
		writeError(w, htsErr)
		return
	}

	var fields []string
	if !strings.HasPrefix(id, "10X") {
		fields, err = parseFields(params)
		if err != nil {
			htsErr := &htsgetError{
				Code: http.StatusBadRequest,
				Htsget: errorContainer{
					"InvalidInput",
					"The request parameters do not adhere to the specification",
				},
			}
			writeError(w, htsErr)
			return
		}
	}

	tags := parseTags(params)

	notags, err := parseNoTags(params, tags)
	if err != nil {
		htsErr := &htsgetError{
			Code: http.StatusBadRequest,
			Htsget: errorContainer{
				"InvalidInput",
				"The request parameters do not adhere to the specification",
			},
		}
		writeError(w, htsErr)
		return
	}

	region := &genomics.Region{Name: refName, Start: start, End: end}

	if refName != "" && refName != "*" {
		if ok, err := referenceExists(id, refName); !ok {
			htsErr := &htsgetError{
				Code: http.StatusNotFound,
				Htsget: errorContainer{
					"NotFound",
					"The resource requested was not found",
				},
			}
			writeError(w, htsErr)
			return
		} else if err != nil {
			writeError(w, err)
		}
	}

	numBytes := res.ContentLength
	var numBlocks int
	var blockSize int64 = 1e9
	if refName != "" {
		numBlocks = 1
	} else {
		if len(fields) == 0 {
			numBlocks = int(math.Ceil(float64(numBytes) / float64(blockSize)))
		}
	}

	// build HTTP response
	u := make([]urlJSON, 0)
	var h *headers
	dataEndpoint, err := getDataURL(region, fields, tags, notags, id, queryClass, host)
	if err != nil {
		writeError(w, err)
	}

	if queryClass == "header" {
		h = &headers{
			BlockID:   "1",
			NumBlocks: "1",
			Class:     "header",
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, "header"})
	} else if len(fields) == 0 && len(tags) == 0 && len(notags) == 0 && refName == "" {
		path := dataSource + filePath(id)
		var start, end int64 = 0, 0

		for i := 1; i <= numBlocks; i++ {
			end = start + blockSize - 1
			if end >= numBytes {
				end = numBytes - 1
			}
			h := &headers{
				Range: "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10),
			}
			start = end + 1
			u = append(u, urlJSON{path, h, ""})
		}
	} else {
		h = &headers{
			BlockID:   "1",
			NumBlocks: "2",
			Class:     "header",
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, "header"})

		h = &headers{
			BlockID:   "2",
			NumBlocks: "2",
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, "body"})
	}

	c := container{format, u, ""}
	ticket := ticket{HTSget: c}
	ct := "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
	w.Header().Set("Content-Type", ct)
	json.NewEncoder(w).Encode(ticket)
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

func getDataURL(r *genomics.Region, fields, tags, notags []string, id, class, host string) (*url.URL, error) {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse(host + "data/")
	if err != nil {
		return nil, err
	}

	// add id url param
	dataEndpoint.Path += id

	// add query params
	query := dataEndpoint.Query()
	if class == "header" {
		query.Set("class", class)
	}
	if r != nil {
		if r.Name != "" {
			query.Set("referenceName", r.Name)
		}
		if r.Start != "-1" {
			query.Set("start", r.Start)
		}
		if r.End != "-1" {
			query.Set("end", r.End)
		}
	}

	if f := strings.Join(fields, ","); f != "" {
		query.Set("fields", f)
	}

	// t := strings.Join(tags, ",")
	// query.Set("tags", t)

	if nt := strings.Join(notags, ","); nt != "" {
		query.Set("notags", nt)
	}

	dataEndpoint.RawQuery = query.Encode()

	return dataEndpoint, nil
}

func referenceExists(id string, refName string) (bool, error) {
	cmd := exec.Command("samtools", "view", "-H", dataSource+filePath(id))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	cmd.Start()
	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()

	for ; err == nil; l, _, err = reader.ReadLine() {
		if strings.Contains(string(l), "SN:"+refName) {
			return true, nil
		}
	}
	cmd.Wait()
	return false, nil
}

func writeError(w http.ResponseWriter, err error) {
	if err, ok := err.(*htsgetError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"htsget": err.Htsget,
		})
		return
	}

	http.Error(w, err.Error(), 500)
}
