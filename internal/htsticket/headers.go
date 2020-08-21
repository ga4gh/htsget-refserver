// Package htsticket ...
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

func NewHeaders() *Headers {
	return new(Headers)
}

func (headers *Headers) SetBlockID(blockID string) *Headers {
	headers.BlockID = blockID
	return headers
}

func (headers *Headers) SetNumBlocks(numBlocks string) *Headers {
	headers.NumBlocks = numBlocks
	return headers
}

func (headers *Headers) SetRangeHeader(start int64, end int64) *Headers {
	headers.Range = "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
	return headers
}

func (headers *Headers) setClass(class string) *Headers {
	headers.Class = class
	return headers
}

func (headers *Headers) SetClassHeader() *Headers {
	headers.setClass(htsconstants.ClassHeader)
	return headers
}

func (headers *Headers) SetClassBody() *Headers {
	headers.setClass(htsconstants.ClassBody)
	return headers
}

func (headers *Headers) SetFilePathHeader(filePath string) *Headers {
	headers.FilePath = filePath
	return headers
}
