package htsserver

import (
	"bufio"
	"io"
	"net/http"
	"os/exec"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htserror"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getVariantsData(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointVariantsData,
		noAfterSetup,
		getVariantsDataHandler,
	).handleRequest(writer, request)
}

// getVariantsData serves the actual data from AWS back to client
func getVariantsDataHandler(handler *requestHandler) {

	fileURL, err := htsconfig.GetObjectPath(handler.HtsReq.GetEndpoint(), handler.HtsReq.GetID())
	if err != nil {
		return
	}

	command, args := constructBcftoolsCommand(handler.HtsReq, fileURL)
	cmd := exec.Command(command, args...)
	pipe, err := cmd.StdoutPipe()

	if err != nil {
		msg := err.Error()
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	err = cmd.Start()
	if err != nil {
		msg := err.Error()
		htserror.InternalServerError(handler.Writer, &msg)
		return
	}

	reader := bufio.NewReader(pipe)
	io.Copy(handler.Writer, reader)
	cmd.Wait()
}

func constructBcftoolsCommand(htsgetReq *htsrequest.HtsgetRequest, fileURL string) (string, []string) {
	command := "bcftools"
	args := []string{"view", fileURL}

	// translate "format" param into bcftools command
	args = append(args, "-O", "v")      // request uncompressed VCF
	args = append(args, "--no-version") // do not add bcftools version to VCF header

	// translate "HtsgetBlockClass" param into bcftools command
	if htsgetReq.GetHtsgetBlockClass() == "header" {
		args = append(args, "-h")
	} else {
		args = append(args, "-H")

		// translate "referenceName", "start", "end" params into bcftools command
		start := htsgetReq.GetStart()
		end := htsgetReq.GetEnd()
		if htsgetReq.ReferenceNameRequested() {
			region := &htsrequest.Region{
				ReferenceName: htsgetReq.GetReferenceName(),
				Start:         &start,
				End:           &end,
			}
			args = append(args, "-r", region.ExportBcftools())
		}
	}
	return command, args
}
