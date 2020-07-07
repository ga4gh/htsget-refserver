package htsgetserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig/htsgetconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsgeterror"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetticket"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgetdao"
)

func getReadsTicket(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	host := htsgetconfig.GetHost()
	htsgetReq, err := htsgetrequest.ReadsTicketEndpointSetAllParameters(request, writer, params)

	if err != nil {
		return
	}

	dao, err := htsgetdao.GetReadsDaoForID(htsgetReq.ID())

	if err != nil {
		msg := "Could not determine data source url from request id"
		htsgeterror.InternalServerError(writer, &msg)
		return
	}

	// build HTTP response
	var urls []*htsgetticket.URL
	dataEndpoint, err := getDataURL(htsgetReq, host)
	if err != nil {
		msg := "Could not construct data url"
		htsgeterror.InternalServerError(writer, &msg)
	}

	if htsgetReq.HeaderOnlyRequested() {
		headers := htsgetticket.NewHeaders().SetBlockID("1").SetNumBlocks("1").SetClassHeader()
		url := htsgetticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headers).SetClassHeader()
		urls = append(urls, url)
	} else if htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested() && htsgetReq.AllRegionsRequested() {
		urls = dao.GetByteRangeUrls()
	} else {
		headersBlock1 := htsgetticket.NewHeaders().SetBlockID("1").SetNumBlocks("2").SetClassHeader()
		urlBlock1 := htsgetticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock1).SetClassHeader()
		headersBlock2 := htsgetticket.NewHeaders().SetBlockID("2").SetNumBlocks("2")
		urlBlock2 := htsgetticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock2)
		urls = append(urls, urlBlock1)
		urls = append(urls, urlBlock2)
	}

	container := htsgetticket.NewContainer().SetFormatBam().SetURLS(urls)
	ticket := htsgetticket.NewTicket().SetContainer(container)
	ct := "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
	writer.Header().Set("Content-Type", ct)
	json.NewEncoder(writer).Encode(ticket)
}

func getDataURL(htsgetReq *htsgetrequest.HtsgetRequest, host string) (*url.URL, error) {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse(host + htsgetconstants.ReadsDataURLPath)
	if err != nil {
		return nil, err
	}
	// add id url param
	dataEndpoint.Path += htsgetReq.ID()
	// add query params
	query := dataEndpoint.Query()
	if htsgetReq.HeaderOnlyRequested() {
		query.Set("class", htsgetReq.Class())
	}
	if htsgetReq.ReferenceNameRequested() {
		query.Set("referenceName", htsgetReq.ReferenceName())
	}
	if htsgetReq.StartRequested() {
		query.Set("start", htsgetReq.Start())
	}
	if htsgetReq.EndRequested() {
		query.Set("end", htsgetReq.End())
	}
	if !htsgetReq.AllFieldsRequested() {
		f := strings.Join(htsgetReq.Fields(), ",")
		query.Set("fields", f)
	}
	if !htsgetReq.TagsNotSpecified() {
		t := strings.Join(htsgetReq.Tags(), ",")
		query.Set("tags", t)
	}
	if !htsgetReq.NoTagsNotSpecified() {
		nt := strings.Join(htsgetReq.NoTags(), ",")
		query.Set("notags", nt)
	}
	dataEndpoint.RawQuery = query.Encode()
	return dataEndpoint, nil
}
