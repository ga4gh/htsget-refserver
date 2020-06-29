// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

// default values for scalar params if param is not specified in request
var defaultScalarParameterValues = map[string]string{
	"id":            "",
	"format":        "BAM",
	"class":         "",
	"referenceName": "*",
	"start":         "-1",
	"end":           "-1",
}

// default values for list params if param is not specified in request
var defaultListParameterValues = map[string][]string{
	"fields": {"ALL"},
	"tags":   {"ALL"},
	"notags": {"NONE"},
}
