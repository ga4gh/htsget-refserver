package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htscli"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getVariantsData(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointVariantsData,
		addRegionFromQueryString,
		getVariantsDataHandler,
	).handleRequest(writer, request)
}

// getVariantsData serves the actual data from AWS back to client
func getVariantsDataHandler(handler *requestHandler) {
	fileURL, err := htsconfig.GetObjectPath(handler.HtsReq.GetEndpoint(), handler.HtsReq.GetID())
	if err != nil {
		return
	}

	commandChain := htscli.NewCommandChain()
	removedHeadBytes := 0
	removedTailBytes := 0

	if handler.HtsReq.IsHeaderBlock() {
		// only get the header for header blocks
		commandChain.AddCommand(bcftoolsViewHeaderOnlyVCF(fileURL))
	} else {
		// body-based requests
		commandChain.AddCommand(bcftoolsViewBodyVCF(handler.HtsReq, fileURL))
	}

	// execute command chain and stream output
	commandWriteStream(commandChain, removedHeadBytes, removedTailBytes, handler.Writer)
}

func bcftoolsViewHeaderOnlyVCF(fileURL string) *htscli.Command {
	cmd := htscli.BcftoolsView()
	cmd.SetFilePath(fileURL)
	cmd.SetHeaderOnly(true)
	return cmd.GetCommand()
}

func bcftoolsViewBodyVCF(htsgetReq *htsrequest.HtsgetRequest, fileURL string) *htscli.Command {
	cmd := htscli.BcftoolsView()
	cmd.SetFilePath(fileURL)
	cmd.SetHeaderOnly(false)
	if !htsgetReq.AllRegionsRequested() {
		cmd.SetRegion(htsgetReq.GetRegions()[0])
	}
	return cmd.GetCommand()
}
