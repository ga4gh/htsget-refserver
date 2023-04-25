// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module parsing contains operations for parsing various parameter types
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

	log "github.com/sirupsen/logrus"
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
		log.Debug("error in params length")
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

// partialRequestBody interface enabling the return of a single attribute from
// the request body
type partialRequestBody interface {
	getattr() reflect.Value
}

// partialRequestBodyFormat parses single format parameter from request body
type partialRequestBodyFormat struct {
	Format *string `json:"format"`
}

// partialRequestBodyFields parses single fields parameter from request body
type partialRequestBodyFields struct {
	Fields *[]string `json:"fields"`
}

// partialRequestBodyTags parses single tags parameter from request body
type partialRequestBodyTags struct {
	Tags *[]string `json:"tags"`
}

// partialRequestBodyNoTags parses single notags parameter from request body
type partialRequestBodyNoTags struct {
	NoTags *[]string `json:"notags"`
}

// partialRequestBodyRegions parses single regions parameter from request body
type partialRequestBodyRegions struct {
	Regions *[]*Region `json:"regions"`
}

// getattr returns reflected format value
func (rb *partialRequestBodyFormat) getattr() reflect.Value {
	return reflect.ValueOf(rb.Format)
}

// getattr returns reflected fields value
func (rb *partialRequestBodyFields) getattr() reflect.Value {
	return reflect.ValueOf(rb.Fields)
}

// getattr returns reflected tags value
func (rb *partialRequestBodyTags) getattr() reflect.Value {
	return reflect.ValueOf(rb.Tags)
}

// getattr returns reflected notags value
func (rb *partialRequestBodyNoTags) getattr() reflect.Value {
	return reflect.ValueOf(rb.NoTags)
}

// getattr returns reflected regions value
func (rb *partialRequestBodyRegions) getattr() reflect.Value {
	return reflect.ValueOf(rb.Regions)
}

// newPartialRequestBody constructs an empty partialRequestBody, holding a
// single parameter based on the passed key
func newPartialRequestBody(key string) partialRequestBody {

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

// parseReqBodyParam parses a single parameter from an overall request body
func parseReqBodyParam(requestBodyBytes []byte, key string) (reflect.Value, bool, error) {

	// construct a single parameter, partial request body and unmarshal JSON
	partialRequestBodyObj := newPartialRequestBody(key)
	err := json.Unmarshal(requestBodyBytes, partialRequestBodyObj)
	if err != nil {
		log.Debugf("error unmarshaling in parseReqBodyParam, %v", err)
		msg := "Could not parse request body, offending attribute: '" + key + "'. Value is malformed or incorrect datatype"
		return reflect.ValueOf(nil), false, errors.New(msg)
	}

	// checks if the value is nil (ie. a nil pointer means no value was found)
	reflectedPtr := partialRequestBodyObj.getattr()
	reflectedValue := reflectedPtr.Elem()
	found := !reflectedPtr.IsNil()
	return reflectedValue, found, nil
}
