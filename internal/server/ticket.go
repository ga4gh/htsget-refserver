package server

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgeterror"

	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/genomics"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgethttp/htsgetrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgetparameters"
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

// Headers contains any headers needed by the server from the client
type headers struct {
	BlockID   string `json:"block-id,omitempty"`   // id of current block
	NumBlocks string `json:"num-blocks,omitempty"` // total number of blocks
	Range     string `json:"range,omitempty"`
	Class     string `json:"class,omitempty"`
}

func getReads(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	host := htsgetutils.AddTrailingSlash(config.GetConfigProp("host"))
	htsgetReq, err := htsgetparameters.ReadsEndpointSetAllParameters(request, writer, params)

	if err != nil {
		return
	}

	region := &genomics.Region{
		Name:  htsgetReq.GetScalar("referenceName"),
		Start: htsgetReq.GetScalar("start"),
		End:   htsgetReq.GetScalar("end"),
	}
	res, _ := http.Head(config.DATA_SOURCE_URL + htsgetutils.FilePath(htsgetReq.GetScalar("id")))
	numBytes := res.ContentLength
	var numBlocks int
	var blockSize int64 = 1e9
	if htsgetReq.GetScalar("referenceName") != "" {
		numBlocks = 1
	} else {
		if len(htsgetReq.GetList("fields")) == 0 {
			numBlocks = int(math.Ceil(float64(numBytes) / float64(blockSize)))
		}
	}

	// build HTTP response
	u := make([]urlJSON, 0)
	var h *headers
	dataEndpoint, err := getDataURL(region, htsgetReq, host)
	if err != nil {
		msg := "Could not construct data url"
		htsgeterror.InternalServerError(writer, &msg)
	}

	if htsgetReq.GetScalar("class") == "header" {
		h = &headers{
			BlockID:   "1",
			NumBlocks: "1",
			Class:     "header",
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, "header"})
	} else if len(htsgetReq.GetList("fields")) == 0 && len(htsgetReq.GetList("tags")) == 0 && len(htsgetReq.GetList("notags")) == 0 && htsgetReq.GetScalar("referenceName") == "*" {
		path := config.DATA_SOURCE_URL + htsgetutils.FilePath(htsgetReq.GetScalar("id"))
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

	c := container{htsgetReq.GetScalar("format"), u, ""}
	ticket := ticket{HTSget: c}
	ct := "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
	writer.Header().Set("Content-Type", ct)
	json.NewEncoder(writer).Encode(ticket)
}

func getDataURL(r *genomics.Region, htsgetReq *htsgetrequest.HtsgetRequest, host string) (*url.URL, error) {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse(host + "data/")
	if err != nil {
		return nil, err
	}

	// add id url param
	dataEndpoint.Path += htsgetReq.GetScalar("id")

	// add query params
	query := dataEndpoint.Query()
	if htsgetReq.GetScalar("class") == "header" {
		query.Set("class", htsgetReq.GetScalar("class"))
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

	if f := strings.Join(htsgetReq.GetList("fields"), ","); f != "" {
		query.Set("fields", f)
	}

	// t := strings.Join(tags, ",")
	// query.Set("tags", t)

	if nt := strings.Join(htsgetReq.GetList("notags"), ","); nt != "" {
		query.Set("notags", nt)
	}

	dataEndpoint.RawQuery = query.Encode()

	return dataEndpoint, nil
}
