package htsserver

import (
	"bufio"
	"io"
	"net/http"
	"os/exec"

	"github.com/ga4gh/htsget-refserver/internal/htscli"
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getReadsData(writer http.ResponseWriter, request *http.Request) {
	newRequestHandler(
		htsconstants.GetMethod,
		htsconstants.APIEndpointReadsData,
		addRegionFromQueryString,
		getReadsDataHandler,
	).handleRequest(writer, request)
}

func getReadsDataHandler(handler *requestHandler) {
	fileURL, err := htsconfig.GetObjectPath(handler.HtsReq.GetEndpoint(), handler.HtsReq.GetID())
	if err != nil {
		return
	}

	commandChain := htscli.NewCommandChain()
	removedHeadBytes := 0
	removedTailBytes := htsconstants.BamEOFLen

	if handler.HtsReq.IsHeaderBlock() {
		// only get the header for header blocks
		commandChain.AddCommand(samtoolsViewHeaderOnlyBAM(fileURL))
	} else {
		// body-based requests will remove header bytes, as they are streamed
		// in a different block
		headerByteSize, _ := getHeaderByteSize(handler.HtsReq.GetID(), fileURL)
		removedHeadBytes = headerByteSize
		var region *htsrequest.Region = nil
		if !handler.HtsReq.AllRegionsRequested() {
			region = handler.HtsReq.GetRegions()[0]
		}

		if handler.HtsReq.AllFieldsRequested() && handler.HtsReq.AllTagsRequested() {
			// simple streaming of single block without field/tag modification
			commandChain.AddCommand(samtoolsViewHeaderExcludedBAM(fileURL, region))

		} else {
			// specific fields/tags requested, requires chaining of samtools
			// with htsget-refserver-utils modify sam commands
			commandChain.AddCommand(samtoolsViewHeaderIncludedSAM(fileURL, region))
			commandChain.AddCommand(modifySam(handler.HtsReq))
			commandChain.AddCommand(samtoolsViewSamToBamStream())
		}
	}

	// execute command chain and stream output
	commandWriteStream(commandChain, removedHeadBytes, removedTailBytes, handler.Writer)

	// write EOF on the last block
	if handler.HtsReq.IsFinalBlock() {
		writeBamEOF(handler.Writer)
	}
}

func commandWriteStream(commandChain *htscli.CommandChain, removeHeadBytes int, removeTailBytes int, writer http.ResponseWriter) error {

	commandChain.SetupCommandChain()
	pipe := commandChain.ExecuteCommandChain()
	reader := bufio.NewReader(pipe)
	bufferSize := 65536
	firstLoop := true
	eofNotReached := true

	for ok := true; ok; ok = eofNotReached {
		bufferBytes := make([]byte, bufferSize)
		nBytesRead, _ := io.ReadFull(reader, bufferBytes)

		// indicates this is the last loop
		if nBytesRead != bufferSize {
			// remove all unread bytes after EOF,
			// then remove bytes specified by removeTailBytes
			bufferBytes = bufferBytes[:nBytesRead]
			bufferBytes = bufferBytes[:len(bufferBytes)-removeTailBytes]
			eofNotReached = false
		}

		// if first loop, remove bytes specified by removeHeadBytes
		if firstLoop {
			firstLoop = false
			bufferBytes = bufferBytes[removeHeadBytes:]
		}

		writer.Write(bufferBytes)
	}
	return nil
}

func writeBamEOF(writer http.ResponseWriter) {
	writer.Write(htsconstants.BamEOF)
}

// for header requests
func samtoolsViewHeaderOnlyBAM(fileURL string) *htscli.Command {
	return htscli.SamtoolsView().AddFilePath(fileURL).HeaderOnly().OutputBAM().GetCommand()
}

// requests for all fields/tags
func samtoolsViewHeaderExcludedBAM(fileURL string, region *htsrequest.Region) *htscli.Command {
	samtoolsView := htscli.SamtoolsView().AddFilePath(fileURL).OutputBAM()
	if region != nil {
		samtoolsView.AddRegion(region)
	}
	return samtoolsView.GetCommand()
}

// commands used when custom fields/tags are requested
func samtoolsViewHeaderIncludedSAM(fileURL string, region *htsrequest.Region) *htscli.Command {
	samtoolsView := htscli.SamtoolsView().AddFilePath(fileURL).HeaderIncluded()
	if region != nil {
		samtoolsView.AddRegion(region)
	}
	return samtoolsView.GetCommand()
}

func modifySam(htsgetReq *htsrequest.HtsgetRequest) *htscli.Command {
	modifySam := htscli.ModifySam()
	if !htsgetReq.AllFieldsRequested() {
		modifySam.SetFields(htsgetReq.GetFields())
	}
	if !htsgetReq.TagsNotSpecified() {
		tags := htsgetReq.GetTags()
		if len(tags) == 1 && tags[0] == "" {
			modifySam.SetTags([]string{"NONE"})
		} else {
			modifySam.SetTags(tags)
		}
	}
	if !htsgetReq.NoTagsNotSpecified() {
		modifySam.SetNoTags(htsgetReq.GetNoTags())
	}
	return modifySam.GetCommand()
}

func samtoolsViewSamToBamStream() *htscli.Command {
	return htscli.SamtoolsView().OutputBAM().StreamFromStdin().GetCommand()
}

func getHeaderByteSize(id string, fileURL string) (int, error) {
	cmd := exec.Command("samtools", "view", "-H", "-b", fileURL)
	tmpHeader, err := htsconfig.CreateTempfile(id + "_header")
	if err != nil {
		return 0, err
	}

	cmd.Stdout = tmpHeader
	cmd.Run()

	fi, err := tmpHeader.Stat()
	if err != nil {
		return 0, err
	}

	size := fi.Size() - 28
	tmpHeader.Close()
	htsconfig.RemoveTempfile(tmpHeader)
	return int(size), nil
}
