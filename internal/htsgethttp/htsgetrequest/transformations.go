// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

import (
	"strings"
)

// the correct transformation function for each scalar parameter (transforms the
// raw param string into a mature value usable by downstream htsget functions)
var transformationScalarByParam = map[string]func(string) string{
	"id":            noTransform,
	"format":        strings.ToUpper,
	"class":         strings.ToLower,
	"referenceName": noTransform,
	"start":         noTransform,
	"end":           noTransform,
}

// the correct transformation function for each list parameter
var transformationListByParam = map[string]func(string) []string{
	"fields": splitAndUppercase,
	"tags":   splitOnComma,
	"notags": splitOnComma,
}

// performs no transformation on a request parameter
func noTransform(s string) string {
	return s
}

// splits a single string into a list of strings, using the comma as delimiter
func splitOnComma(s string) []string {
	return strings.Split(s, ",")
}

// splits into a list of strings, and makes each string uppercase
func splitAndUppercase(s string) []string {
	sList := splitOnComma(s)
	for i := 0; i < len(sList); i++ {
		sList[i] = strings.ToUpper(sList[i])
	}
	return sList
}
