// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module defaults.go contains default values for each parameter
package htsrequest

// defaultScalarParameterValues (map[string]string): values for scalar params
// if param is not specified in request
var defaultScalarParameterValues = map[string]string{
	"id":            "",
	"format":        "BAM",
	"class":         "",
	"referenceName": "",
	"start":         "-1",
	"end":           "-1",
}

// defaultListParameterValues (map[string][]string): values for list params
// if param is not specified in request
var defaultListParameterValues = map[string][]string{
	"fields": {"ALL"},
	"tags":   {"ALL"},
	"notags": {"NONE"},
}
