// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module transformations defines operations for transforming the raw string
// parsed from the HTTP request into a mature value that is usable by the program
package htsrequest

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ParamTransformer transforms request parameters on query string to expected datatype
type ParamTransformer struct{}

// NewParamTransformer instantiates a new ParamTransformer object
func NewParamTransformer() *ParamTransformer {
	return new(ParamTransformer)
}

// NoTransform performs no param transformation, returning the exact same value
func (t *ParamTransformer) NoTransform(s string) (string, string) {
	return s, ""
}

// TransformStringUppercase transforms a param with lowercase characters to all
// uppercase
func (t *ParamTransformer) TransformStringUppercase(s string) (string, string) {
	return strings.ToUpper(s), ""
}

// TransformStringLowercase transforms a param with uppercase characters to all
// lowercase
func (t *ParamTransformer) TransformStringLowercase(s string) (string, string) {
	return strings.ToLower(s), ""
}

// TransformStringToInt converts a request param to integer datatype
func (t *ParamTransformer) TransformStringToInt(s string) (int, string) {
	msg := ""
	value, err := strconv.Atoi(s)
	if err != nil {
		log.Debugf("error in TransformStringToInt, %v", err)
		msg = fmt.Sprintf("Could not parse value: '%s', integer expected", s)
	}
	return value, msg
}

// TransformSplit splits a string into a list of strings, delimited by comma
func (t *ParamTransformer) TransformSplit(s string) ([]string, string) {
	return strings.Split(s, ","), ""
}

// TransformSplitAndUppercase splits a string into a list of strings, and
// uppercases each element
func (t *ParamTransformer) TransformSplitAndUppercase(s string) ([]string, string) {
	sList, _ := t.TransformSplit(s)
	for i := 0; i < len(sList); i++ {
		sList[i] = strings.ToUpper(sList[i])
	}
	return sList, ""
}
