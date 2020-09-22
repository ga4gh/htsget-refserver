package htsserver

import (
	"fmt"

	"github.com/ga4gh/htsget-refserver/internal/htsdao"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

func ticketRequestHandler(handler *requestHandler) {
	fmt.Println(handler.HtsReq.ID())
	fmt.Println(handler.HtsReq.ReferenceName())
	fmt.Println("***")

	dao, err := htsdao.GetDao(handler.HtsReq)
	if err != nil {
		msg := "Could not determine data source path/url from request id"
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	var urls []*htsticket.URL
	dataEndpoint, err := handler.HtsReq.ConstructDataEndpointURL()
	if err != nil {
		msg := "Could not construct data url"
		htserror.InternalServerError(handler.Writer, &msg)
	}

	if handler.HtsReq.HeaderOnlyRequested() {
		headers := htsticket.NewHeaders().SetBlockID("1").SetNumBlocks("1").SetClassHeader()
		url := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headers).SetClassHeader()
		urls = append(urls, url)
	} else if handler.HtsReq.AllFieldsRequested() && handler.HtsReq.AllTagsRequested() && handler.HtsReq.AllRegionsRequested() {
		urls = dao.GetByteRangeUrls()
	} else {
		headersBlock1 := htsticket.NewHeaders().SetBlockID("1").SetNumBlocks("2").SetClassHeader()
		urlBlock1 := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock1).SetClassHeader()
		headersBlock2 := htsticket.NewHeaders().SetBlockID("2").SetNumBlocks("2")
		urlBlock2 := htsticket.NewURL().SetURL(dataEndpoint.String()).SetHeaders(headersBlock2)
		urls = append(urls, urlBlock1)
		urls = append(urls, urlBlock2)
	}

	htsticket.FinalizeTicket(handler.HtsReq.Format(), urls, handler.Writer)
}
