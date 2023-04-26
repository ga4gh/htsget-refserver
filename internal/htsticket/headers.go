// Package htsticket produces the htsget JSON response ticket
//
// Module headers holds header keys/values for a single data download url object
package htsticket

import (
	"strconv"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// Headers contains any headers needed by the server from the client
type Headers struct {
	BlockClass    string `json:"HtsgetBlockClass,omitempty"`
	CurrentBlock  string `json:"HtsgetCurrentBlock,omitempty"` // number of current block
	TotalBlocks   string `json:"HtsgetTotalBlocks,omitempty"`  // total number of blocks
	FilePath      string `json:"HtsgetFilePath,omitempty"`
	Range         string `json:"Range,omitempty"`
	Authorization string `json:"Authorization,omitempty"`
}

// NewHeaders instantiates an empty headers object
func NewHeaders() *Headers {
	return new(Headers)
}

// SetCurrentBlock assigns the CurrentBlock header value
func (headers *Headers) SetCurrentBlock(currentBlock string) *Headers {
	headers.CurrentBlock = currentBlock
	return headers
}

// SetTotalBlocks assigns the TotalBlocks header value
func (headers *Headers) SetTotalBlocks(totalBlocks string) *Headers {
	headers.TotalBlocks = totalBlocks
	return headers
}

// SetRangeHeader assigns the Range header value
func (headers *Headers) SetRangeHeader(start int64, end int64) *Headers {
	headers.Range = "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
	return headers
}

// setBlockClass assigns the BlockClass header value
func (headers *Headers) setBlockClass(blockClass string) *Headers {
	headers.BlockClass = blockClass
	return headers
}

// SetClassHeader assigns the Class header value to "header", indicating the data
// download url is responsible for downloading the file header
func (headers *Headers) SetClassHeader() *Headers {
	headers.setBlockClass(htsconstants.ClassHeader)
	return headers
}

// SetClassBody assigns the Class header value to "body", indicating the data
// download url is responsible for downloading the file body
func (headers *Headers) SetClassBody() *Headers {
	headers.setBlockClass(htsconstants.ClassBody)
	return headers
}

// SetFilePathHeader assigns the FilePath header, informing the data or
// file bytes endpoint of which file to stream back to client
func (headers *Headers) SetFilePathHeader(filePath string) *Headers {
	headers.FilePath = filePath
	return headers
}
