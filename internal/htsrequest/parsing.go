// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module parsing.go contains operations for parsing various parameter types
// from an HTTP request
package htsrequest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"
	"github.com/go-chi/chi"
)

// parsePathParam parses a single url path parameter as a string
func parsePathParam(request *http.Request, key string) (string, bool) {
	value := chi.URLParam(request, key)
	found := false
	if !htsutils.StringIsEmpty(value) {
		found = true
	}
	return value, found
}

// parseQueryParam parses a single query string parameter as a string
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
func parseHeaderParam(request *http.Request, key string) (string, bool) {

	value := request.Header.Get(key)
	found := false
	if !htsutils.StringIsEmpty(value) {
		found = true
	}
	return value, found
}

// parsing of partial request body

type partialRequestBody interface {
	getattr() reflect.Value
}

type partialRequestBodyFormat struct {
	Format *string `json:"format"`
}

type partialRequestBodyFields struct {
	Fields *[]string `json:"fields"`
}

type partialRequestBodyTags struct {
	Tags *[]string `json:"tags"`
}

type partialRequestBodyNoTags struct {
	NoTags *[]string `json:"notags"`
}

type partialRequestBodyRegions struct {
	Regions *[]*Region `json:"regions"`
}

func (rb *partialRequestBodyFormat) getattr() reflect.Value {
	return reflect.ValueOf(rb.Format)
}

func (rb *partialRequestBodyFields) getattr() reflect.Value {
	return reflect.ValueOf(rb.Fields)
}

func (rb *partialRequestBodyTags) getattr() reflect.Value {
	return reflect.ValueOf(rb.Tags)
}

func (rb *partialRequestBodyNoTags) getattr() reflect.Value {
	return reflect.ValueOf(rb.NoTags)
}

func (rb *partialRequestBodyRegions) getattr() reflect.Value {
	return reflect.ValueOf(rb.Regions)
}

func NewPartialRequestBody(key string) partialRequestBody {

	var prb partialRequestBody
	switch key {
	case "format":
		prb = new(partialRequestBodyFormat)
	case "fields":
		prb = new(partialRequestBodyFields)
	case "tags":
		prb = new(partialRequestBodyTags)
	case "notags":
		prb = new(partialRequestBodyNoTags)
	case "regions":
		prb = new(partialRequestBodyRegions)
	}
	return prb
}

func parseReqBodyParam(requestBodyBytes []byte, key string) (reflect.Value, bool, error) {

	partialRequestBodyObj := NewPartialRequestBody(key)
	err := json.Unmarshal(requestBodyBytes, partialRequestBodyObj)
	if err != nil {
		msg := "Could not parse request body, offending attribute: '" + key + "'. Value is malformed or incorrect datatype"
		return reflect.ValueOf(nil), false, errors.New(msg)
	}
	reflectedPtr := partialRequestBodyObj.getattr()
	reflectedValue := reflectedPtr.Elem()
	found := !reflectedPtr.IsNil()
	return reflectedValue, found, nil
}
