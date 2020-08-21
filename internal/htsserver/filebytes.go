package htsserver

import (
	"bufio"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

func getFileBytes(writer http.ResponseWriter, request *http.Request) {

	params := request.URL.Query()
	htsgetReq, err := htsrequest.FileBytesEndpointSetAllParameters(request, writer, params)
	if err != nil {
		return
	}

	start, end, err := htsutils.ParseRangeHeader(htsgetReq.Range())
	if err != nil {
		return
	}

	file, err := os.Open(htsgetReq.HtsgetFilePath())
	reader := bufio.NewReader(file)
	reader.Discard(int(start))

	chunkSize := 8192
	nBytes := end - start + 1
	nChunks := int(math.Ceil(float64(nBytes) / float64(chunkSize)))
	nBytesRead := int64(0)
	buffer := make([]byte, chunkSize)

	for i := 0; i < nChunks; i++ {

		n, err := io.ReadFull(reader, buffer)
		if err != nil {

		}
		writer.Write(buffer)
		// at the second last chunk, create a buffer that will only read bytes
		// up to the exact remaining needed to fulfill the range request
		nBytesRead += int64(n)
		if i == nChunks-2 {
			bytesRemaining := nBytes - nBytesRead
			buffer = make([]byte, bytesRemaining)
		} else {
			buffer = make([]byte, chunkSize)
		}
	}
}
