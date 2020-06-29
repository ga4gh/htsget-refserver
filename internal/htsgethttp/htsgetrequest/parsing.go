// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

// indicates whether each htsget parameter is found on the url path,
// if false, it is found in the query string
var isPathByParam = map[string]bool{
	"id":            true,
	"format":        false,
	"class":         false,
	"referenceName": false,
	"start":         false,
	"end":           false,
	"fields":        false,
	"tags":          false,
	"notags":        false,
}

// indicates whether each htsget parameter is expected to contain a scalar
// value, if false, it contains a list value
var isScalarByParam = map[string]bool{
	"id":            true,
	"format":        true,
	"class":         true,
	"referenceName": true,
	"start":         true,
	"end":           true,
	"fields":        false,
	"tags":          false,
	"notags":        false,
}

// parse a single url path parameter as a string
func parsePathParam(request *http.Request, key string) string {
	value := chi.URLParam(request, key)
	return value
}

// parse a single query string parameter as a string
func parseQueryParam(params url.Values, key string) (string, error) {
	if len(params[key]) == 1 {
		return params[key][0], nil
	}
	if len(params[key]) > 1 {
		return "", errors.New("too many values specified for parameter: " + key)
	}
	return "", nil
}
