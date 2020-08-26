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
	BlockID   string `json:"HtsgetBlockId,omitempty"`   // id of current block
	NumBlocks string `json:"HtsgetNumBlocks,omitempty"` // total number of blocks
	Range     string `json:"Range,omitempty"`
	Class     string `json:"HtsgetBlockClass,omitempty"`
	FilePath  string `json:"HtsgetFilePath,omitempty"`
}

// NewHeaders instantiates an empty headers object
func NewHeaders() *Headers {
	return new(Headers)
}

// SetBlockID assigns the BlockID header value
func (headers *Headers) SetBlockID(blockID string) *Headers {
	headers.BlockID = blockID
	return headers
}

// SetNumBlocks assigns the NumBlocks header value
func (headers *Headers) SetNumBlocks(numBlocks string) *Headers {
	headers.NumBlocks = numBlocks
	return headers
}

// SetRangeHeader assigns the Range header value
func (headers *Headers) SetRangeHeader(start int64, end int64) *Headers {
	headers.Range = "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
	return headers
}

// setClass assigns the Class header value
func (headers *Headers) setClass(class string) *Headers {
	headers.Class = class
	return headers
}

// SetClassHeader assigns the Class header value to "header", indicating the data
// download url is responsible for downloading the file header
func (headers *Headers) SetClassHeader() *Headers {
	headers.setClass(htsconstants.ClassHeader)
	return headers
}

// SetClassBody assigns the Class header value to "body", indicating the data
// download url is responsible for downloading the file body
func (headers *Headers) SetClassBody() *Headers {
	headers.setClass(htsconstants.ClassBody)
	return headers
}

// SetFilePathHeader assigns the FilePath header, informing the data or
// file bytes endpoint of which file to stream back to client
func (headers *Headers) SetFilePathHeader(filePath string) *Headers {
	headers.FilePath = filePath
	return headers
}
