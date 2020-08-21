// Package htsticket ...
package htsticket

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// URL holds the url, headers and class
type URL struct {
	URL     string   `json:"url"`
	Headers *Headers `json:"headers,omitempty"`
	Class   string   `json:"class,omitempty"`
}

func NewURL() *URL {
	return new(URL)
}

func (urlObj *URL) SetURL(url string) *URL {
	urlObj.URL = url
	return urlObj
}

func (urlObj *URL) SetHeaders(headers *Headers) *URL {
	urlObj.Headers = headers
	return urlObj
}

func (urlObj *URL) setClass(class string) *URL {
	urlObj.Class = class
	return urlObj
}

func (urlObj *URL) SetClassHeader() *URL {
	urlObj.setClass(htsconstants.ClassHeader)
	return urlObj
}

func (urlObj *URL) SetClassBody() *URL {
	urlObj.setClass(htsconstants.ClassBody)
	return urlObj
}
