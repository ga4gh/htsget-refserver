// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module transformations.go defines operations for transforming the raw string
// parsed from the HTTP request into a mature value that is usable by the program
package htsrequest

import (
	"strings"
)

// transformationScalarByParam (map[string]func(string) string): map of
// functions. the correct transformation function for each scalar parameter
var transformationScalarByParam = map[string]func(string) string{
	"id":               noTransform,
	"format":           strings.ToUpper,
	"class":            strings.ToLower,
	"referenceName":    noTransform,
	"start":            noTransform,
	"end":              noTransform,
	"HtsgetBlockClass": strings.ToLower,
	"HtsgetBlockId":    noTransform,
	"HtsgetNumBlocks":  noTransform,
	"HtsgetFilePath":   noTransform,
	"Range":            noTransform,
}

// transformationScalarByParam (map[string]func(string) []string): map of
// functions. the correct transformation function for each list parameter
var transformationListByParam = map[string]func(string) []string{
	"fields": splitAndUppercase,
	"tags":   splitOnComma,
	"notags": splitOnComma,
}

// noTransform performs no transformation on a request parameter
//
// Arguments
//	s (string): request parameter
// Returns
//	(string): unmodified request parameter
func noTransform(s string) string {
	return s
}

// splitOnComma splits a single string into a list of strings, using the comma
// as delimiter
//
// Arguments
//	s (string): request parameter string
// Returns
//	([]string): list representation of the passed string
func splitOnComma(s string) []string {
	return strings.Split(s, ",")
}

// splitAndUppercase splits into a list of strings, and makes each string
// uppercase
//
// Arguments
//	s (string): request parameter string
// Returns
// ([]string): list of strings, with each string uppercased
func splitAndUppercase(s string) []string {
	sList := splitOnComma(s)
	for i := 0; i < len(sList); i++ {
		sList[i] = strings.ToUpper(sList[i])
	}
	return sList
}
