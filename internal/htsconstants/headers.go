// Package htsconstants contains program constants
//
// Module headers contains constants relating to http header values used by htsget
package htsconstants

/* **************************************************
 * HTTP HEADER NAMES
 * ************************************************** */

// HTTPHeaderName enum for http header names/keys
type HTTPHeaderName int

// enum values for HttpHeaderName
const (
	ContentTypeHeader HTTPHeaderName = 0
)

// string representations of HttpHeaderName enum
const (
	contentTypeHeaderString = "Content-Type"
)

// httpHeaderNameStringMap maps HttpHeaderName enum values to string representation
var httpHeaderNameStringMap = map[HTTPHeaderName]string{
	ContentTypeHeader: contentTypeHeaderString,
}

// String gets the string representation of a HttpHeaderName enum instance
func (e HTTPHeaderName) String() string {
	return httpHeaderNameStringMap[e]
}

/* **************************************************
 * CONTENT-TYPE HEADER VALUES
 * ************************************************** */

// ContentTypeHeaderValue enum for content-type header values
type ContentTypeHeaderValue int

// enum values for ContentTypeHeaderValue
const (
	ContentTypeHeaderHtsgetJSON ContentTypeHeaderValue = 0
)

// string representations of ContentTypeHeaderValue enum
const (
	contentTypeHeaderHtsgetJSONString string = "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
)

// contentTypeStringMap maps ContentTypeHeaderValue enum values to string representation
var contentTypeStringMap = map[ContentTypeHeaderValue]string{
	ContentTypeHeaderHtsgetJSON: contentTypeHeaderHtsgetJSONString,
}

// String gets the string representation of a ContentTypeHeaderValue enum instance
func (e ContentTypeHeaderValue) String() string {
	return contentTypeStringMap[e]
}
