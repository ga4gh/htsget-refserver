// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module parsing.go contains operations for parsing various parameter types
// from an HTTP request
package htsrequest

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"
	"github.com/go-chi/chi"
)

type PostRequestBody struct {
	Format  *string    `json:"format"`
	Fields  *[]string  `json:"fields"`
	Tags    *[]string  `json:"tags"`
	NoTags  *[]string  `json:"notags"`
	Regions *[]*Region `json:"regions"`
}

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

func parseReqBodyParam(requestBody *PostRequestBody, key string) (string, bool) {
	fmt.Println("Your requst body")
	fmt.Println(requestBody)

	foo := reflect.TypeOf(requestBody)
	fmt.Println(foo)

	/*
		requestBodyReflect := reflect.ValueOf(requestBody).Type()
		requestFieldReflect := requestBodyReflect.FieldByName(key)
		fmt.Println(requestFieldReflect)

		fmt.Println("Your key")
		fmt.Println(key)
	*/
	return "", false
}
