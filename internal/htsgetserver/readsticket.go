package htsgetserver

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/genomics"
	"github.com/ga4gh/htsget-refserver/internal/htsgeterror"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
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
	BlockID   string `json:"HtsgetBlockId,omitempty"`   // id of current block
	NumBlocks string `json:"HtsgetNumBlocks,omitempty"` // total number of blocks
	Range     string `json:"Range,omitempty"`
	Class     string `json:"HtsgetBlockClass,omitempty"`
}

func getReadsTicket(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	host := htsgetutils.AddTrailingSlash(config.GetConfigProp("host"))
	htsgetReq, err := htsgetrequest.ReadsTicketEndpointSetAllParameters(request, writer, params)

	if err != nil {
		return
	}

	region := &genomics.Region{
		Name:  htsgetReq.ReferenceName(),
		Start: htsgetReq.Start(),
		End:   htsgetReq.End(),
	}
	res, _ := http.Head(config.DATA_SOURCE_URL + htsgetutils.FilePath(htsgetReq.ID()))
	numBytes := res.ContentLength
	var numBlocks int
	var blockSize int64 = 1e9
	if htsgetReq.ReferenceName() != "" {
		numBlocks = 1
	} else {
		if len(htsgetReq.Fields()) == 0 {
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

	if htsgetReq.Class() == "header" {
		h = &headers{
			BlockID:   "1",
			NumBlocks: "1",
			Class:     "header",
		}
		u = append(u, urlJSON{dataEndpoint.String(), h, "header"})
	} else if htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested() && htsgetReq.ReferenceName() == "*" {
		path := config.DATA_SOURCE_URL + htsgetutils.FilePath(htsgetReq.ID())
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

	c := container{htsgetReq.Format(), u, ""}
	ticket := ticket{HTSget: c}
	ct := "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
	writer.Header().Set("Content-Type", ct)
	json.NewEncoder(writer).Encode(ticket)
}

func getDataURL(r *genomics.Region, htsgetReq *htsgetrequest.HtsgetRequest, host string) (*url.URL, error) {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse(host + "reads/data/")
	if err != nil {
		return nil, err
	}

	// add id url param
	dataEndpoint.Path += htsgetReq.ID()

	// add query params
	query := dataEndpoint.Query()
	if htsgetReq.Class() == "header" {
		query.Set("class", htsgetReq.Class())
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

	if !htsgetReq.AllFieldsRequested() {
		f := strings.Join(htsgetReq.Fields(), ",")
		query.Set("fields", f)
	}

	//if t := strings.Join(htsgetReq.Tags(), ","); t != "" {
	//	query.Set("tags", t)
	//}

	if !htsgetReq.AllTagsRequested() {
		nt := strings.Join(htsgetReq.NoTags(), ",")
		query.Set("notags", nt)
	}

	dataEndpoint.RawQuery = query.Encode()

	return dataEndpoint, nil
}
