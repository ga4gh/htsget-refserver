// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module parsing.go contains operations for parsing various parameter types
// from an HTTP request
package htsrequest

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"
	"github.com/go-chi/chi"
)

// ParamType enum of the final data type of an htsget parameter
type ParamType int

// enum values of ParamType: Scalar, List
const (
	ParamTypeScalar ParamType = iota
	ParamTypeList
)

type PostRequestBody struct {
	Format  string   `json:"format"`
	Fields  []string `json:"fields"`
	Tags    []string `json:tags`
	NoTags  []string `json:notags`
	Regions []Region `json:"regions"`
}

// paramTypes (map[string]ParamType): indicates whether each htsget parameter is
// expected to contain a scalar or list value
var paramTypes = map[string]ParamType{
	"id":               ParamTypeScalar,
	"format":           ParamTypeScalar,
	"class":            ParamTypeScalar,
	"referenceName":    ParamTypeScalar,
	"start":            ParamTypeScalar,
	"end":              ParamTypeScalar,
	"fields":           ParamTypeList,
	"tags":             ParamTypeList,
	"notags":           ParamTypeList,
	"HtsgetBlockClass": ParamTypeScalar,
	"HtsgetBlockId":    ParamTypeScalar,
	"HtsgetNumBlocks":  ParamTypeScalar,
	"HtsgetFilePath":   ParamTypeScalar,
	"Range":            ParamTypeScalar,
}

// parsePathParam parses a single url path parameter as a string
//
// Arguments
//	request (*http.Request): the HTTP request
//	key (string): the parameter name/field
// Returns
//	(string): the value of the path parameter by specified name
//	(bool): true if the parameter was specified by client
func parsePathParam(request *http.Request, key string) (string, bool) {
	value := chi.URLParam(request, key)
	found := false
	if !htsutils.StringIsEmpty(value) {
		found = true
	}
	return value, found
}

// parseQueryParam parses a single query string parameter as a string
//
// Arguments
//	params (url.Values): query string parameters from HTTP request
//	key (string): the parameter name/field
// Returns
//	(string): the value of the query parameter by specified name
//	(bool): true if the parameter was specified by client
//	(error): encountered if the parameter was specified by client incorrectly
func parseQueryParam(params url.Values, key string) (string, bool, error) {

	if len(params[key]) == 1 {
		return params[key][0], true, nil
	}
	if len(params[key]) > 1 {
		return "", true, errors.New("too many values specified for parameter: " + key)
	}
	return "", false, nil
}

// parseHeaderParam parses a single header parameter as a string
//
// Arguments
//	request (*http.Request): the HTTP request
//	key (string): the parameter name/field
// Returns
//	(string): the value of the header parameter by specified name
//	(bool): true if the parameter was specified by client
func parseHeaderParam(request *http.Request, key string) (string, bool) {

	value := request.Header.Get(key)
	found := false
	if !htsutils.StringIsEmpty(value) {
		found = true
	}
	return value, found
}

func parseReqBodyParam(request *http.Request, key string) (string, bool) {
	return "", false
}
