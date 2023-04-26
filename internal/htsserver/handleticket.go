package htsserver

import (
	"strconv"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	log "github.com/sirupsen/logrus"

	"github.com/ga4gh/htsget-refserver/internal/htsdao"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

func addBlockURL(blockURLs []*htsticket.URL, blockURL *htsticket.URL) []*htsticket.URL {
	return append(blockURLs, blockURL)
}

func addHeaderBlockURL(blockURLs []*htsticket.URL, request *htsrequest.HtsgetRequest, totalBlocks int) []*htsticket.URL {
	blockHeaders := htsticket.NewHeaders().
		SetCurrentBlock("0").
		SetTotalBlocks(strconv.Itoa(totalBlocks)).
		SetClassHeader()
	dataEndpoint, _ := request.ConstructDataEndpointURL(false, 0)
	blockURL := htsticket.NewURL().
		SetURL(dataEndpoint).
		SetHeaders(blockHeaders).
		SetClassHeader()
	return addBlockURL(blockURLs, blockURL)
}

func addBodyBlockURL(blockURLs []*htsticket.URL, request *htsrequest.HtsgetRequest, currentBlock int, totalBlocks int, useRegion bool, regionI int) []*htsticket.URL {
	blockHeaders := htsticket.NewHeaders().
		SetCurrentBlock(strconv.Itoa(currentBlock)).
		SetTotalBlocks(strconv.Itoa(totalBlocks))
	dataEndpoint, _ := request.ConstructDataEndpointURL(useRegion, regionI)
	blockURL := htsticket.NewURL().
		SetURL(dataEndpoint).
		SetHeaders(blockHeaders).
		SetClassBody()
	return addBlockURL(blockURLs, blockURL)
}

func ticketRequestHandler(handler *requestHandler) {

	dao, err := htsdao.GetDao(handler.HtsReq)
	if err != nil {
		log.Errorf("Could not determine data source path/url from request id, %v", err)
		msg := "Could not determine data source path/url from request id"
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	var blockURLs []*htsticket.URL

	// only header is requested, requires one URL block
	if handler.HtsReq.HeaderOnlyRequested() {
		blockURLs = addHeaderBlockURL(blockURLs, handler.HtsReq, 1)
		// pure byte range URLs, requires one block per every x bytes
	} else if handler.HtsReq.AllFieldsRequested() && handler.HtsReq.AllTagsRequested() && handler.HtsReq.AllRegionsRequested() {
		blockURLs = dao.GetByteRangeUrls(handler.Request)
	} else {
		if handler.HtsReq.AllRegionsRequested() {
			// the entire file was requested, requires 2 blocks: one for header
			// and one for body
			blockURLs = addHeaderBlockURL(blockURLs, handler.HtsReq, 2)
			blockURLs = addBodyBlockURL(blockURLs, handler.HtsReq, 1, 2, false, 0)
		} else {
			// one or more regions requested, requires one header block, and one
			// block for per region
			nBlocks := handler.HtsReq.NRegions() + 1
			blockURLs = addHeaderBlockURL(blockURLs, handler.HtsReq, nBlocks)
			for i := range handler.HtsReq.GetRegions() {
				blockURLs = addBodyBlockURL(blockURLs, handler.HtsReq, i+1, nBlocks, true, i)
			}
		}
	}
	htsticket.FinalizeTicket(handler.HtsReq.GetFormat(), blockURLs, handler.Writer)
}
