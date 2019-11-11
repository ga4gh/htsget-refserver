package server

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

var EOF, _ = hex.DecodeString("1f8b08040000000000ff0600424302001b0003000000000000000000")

var EOF_LEN = int64(len(EOF))

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
	BlockID   string `json:"block-id"`   // id of current block
	NumBlocks string `json:"num-blocks"` // total number of blocks
	Range     string `json:"range,omitempty"`
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
	queryClass, err := parseQueryClass(params)
	refName, err := parseRefName(params)
	start, end, err := parseRange(params, refName)
	fields, err := parseFields(params)
	if err != nil {
		panic(err)
	}

	if refName != "" && refName != "*" {
		if !referenceExists(id, refName) {
			panic("requested reference does not exist")
		}
	}

	numBlocks := 0
	if numBlocks == 0 {
		numBlocks = 1
	}

	// build HTTP response
	u := make([]urlJSON, 0)
	var respClass string
	if queryClass == "header" {
		respClass = "header"
	}
	dataEndpoint := getDataURL(format, refName, start, end, fields, id)
	var h *headers
	if numBlocks == 1 {
		h = &headers{
			BlockID:   "1",
			NumBlocks: strconv.Itoa(numBlocks),
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, respClass})
	} else {
		//var start int64 = 0
		//var blockSize int64 = int64(math.Floor((float64(numBytes) / float64(numBlocks))))
		//var end int64 = start + blockSize
		for i := 1; i <= numBlocks; i++ {
			if i == 1 { // first of multiple blocks
				h = &headers{
					BlockID:   strconv.Itoa(i),
					NumBlocks: strconv.Itoa(numBlocks),
					//Range:     "bytes=" + "0" + "-" + strconv.FormatInt(hLen, 10),
				}
				//start = hLen + 1
			}
			/* else {*/
			//[>       if end > numBytes {<]
			////end = numBytes
			//[>}<]
			//h = &headers{"bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)}
			/*}*/
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

func getDataURL(format string, refName string, start string, end string, fields []string, id string) *url.URL {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse("localhost:3000/data/")
	if err != nil {
		panic(err)
	}

	// add id url param
	if os.Getenv("APP_ENV") == "production" {
		dataEndpoint.Path += id
	} else {
		dataEndpoint.Opaque += id
	}

	// add query params
	query := dataEndpoint.Query()
	query.Set("format", format)
	if refName != "" {
		query.Set("referenceName", refName)
	}
	if start != "-1" {
		query.Set("start", start)
	}
	if end != "-1" {
		query.Set("end", end)
	}
	if f := strings.Join(fields, ","); f != "" {
		query.Set("fields", f)
	}
	dataEndpoint.RawQuery = query.Encode()

	return dataEndpoint
}

func referenceExists(id string, refName string) bool {
	cmd := exec.Command("samtools", "view", "-H", dataSource+filePath(id))
	pipe, _ := cmd.StdoutPipe()
	cmd.Start()
	reader := bufio.NewReader(pipe)
	l, _, err := reader.ReadLine()

	for ; err == nil; l, _, err = reader.ReadLine() {
		if strings.Contains(string(l), "SN:"+refName) {
			return true
		}
	}
	cmd.Wait()
	return false
}
