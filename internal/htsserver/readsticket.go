package htsserver

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsdao"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

func getReadsTicket(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	host := htsconfig.GetHost()
	htsgetReq, err := htsrequest.ReadsTicketEndpointSetAllParameters(request, writer, params)

	if err != nil {
		return
	}

	dao, err := htsdao.GetReadsDaoForID(htsgetReq.ID())

	if err != nil {
		msg := "Could not determine data source path/url from request id"
		htserror.InternalServerError(writer, &msg)
		return
	}

	// build HTTP response
	var urls []*htsticket.URL
	dataEndpoint, err := getDataURL(htsgetReq, host)
	if err != nil {
		msg := "Could not construct data url"
		htserror.InternalServerError(writer, &msg)
	}

	if htsgetReq.HeaderOnlyRequested() {
		headers := htsticket.NewHeaders().SetBlockID("1").SetNumBlocks("1").SetClassHeader()
		url := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headers).SetClassHeader()
		urls = append(urls, url)
	} else if htsgetReq.AllFieldsRequested() && htsgetReq.AllTagsRequested() && htsgetReq.AllRegionsRequested() {
		urls = dao.GetByteRangeUrls()
	} else {
		headersBlock1 := htsticket.NewHeaders().SetBlockID("1").SetNumBlocks("2").SetClassHeader()
		urlBlock1 := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock1).SetClassHeader()
		headersBlock2 := htsticket.NewHeaders().SetBlockID("2").SetNumBlocks("2")
		urlBlock2 := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock2)
		urls = append(urls, urlBlock1)
		urls = append(urls, urlBlock2)
	}

	container := htsticket.NewContainer().SetFormatBam().SetURLS(urls)
	ticket := htsticket.NewTicket().SetContainer(container)
	ct := "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
	writer.Header().Set("Content-Type", ct)
	json.NewEncoder(writer).Encode(ticket)
}

func getDataURL(htsgetReq *htsrequest.HtsgetRequest, host string) (*url.URL, error) {
	// The address of the endpoint on this server which serves the data
	var dataEndpoint, err = url.Parse(host + htsconstants.ReadsDataURLPath)
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
