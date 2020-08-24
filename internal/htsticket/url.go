// Package htsticket produces the htsget JSON response ticket
//
// Module url holds information for downloading a single filepart from a ticket url
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

// NewURL instantiates an empty url object
func NewURL() *URL {
	return new(URL)
}

// SetURL assign the data download url for the filepart
func (urlObj *URL) SetURL(url string) *URL {
	urlObj.URL = url
	return urlObj
}

// SetHeaders assign all headers necessary to access the data
func (urlObj *URL) SetHeaders(headers *Headers) *URL {
	urlObj.Headers = headers
	return urlObj
}

// setClass assigns the value of the class attribute
func (urlObj *URL) setClass(class string) *URL {
	urlObj.Class = class
	return urlObj
}

// SetClassHeader sets class to "header", indicating the url is responsible for
// downloading the requested file's header
func (urlObj *URL) SetClassHeader() *URL {
	urlObj.setClass(htsconstants.ClassHeader)
	return urlObj
}

// SetClassBody sets class to "body", indicating the url is responsible for
// downloading the requested file's body
func (urlObj *URL) SetClassBody() *URL {
	urlObj.setClass(htsconstants.ClassBody)
	return urlObj
}
