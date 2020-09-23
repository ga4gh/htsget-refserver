// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module transformations.go defines operations for transforming the raw string
// parsed from the HTTP request into a mature value that is usable by the program
package htsrequest

import (
	"fmt"
	"strconv"
	"strings"
)

type ParamTransformer struct{}

func NewParamTransformer() *ParamTransformer {
	return new(ParamTransformer)
}

func (t *ParamTransformer) NoTransform(s string) (string, string) {
	return s, ""
}

func (t *ParamTransformer) TransformStringUppercase(s string) (string, string) {
	return strings.ToUpper(s), ""
}

func (t *ParamTransformer) TransformStringLowercase(s string) (string, string) {
	return strings.ToLower(s), ""
}

func (t *ParamTransformer) TransformStringToInt(s string) (int, string) {
	msg := ""
	value, err := strconv.Atoi(s)
	if err != nil {
		msg = fmt.Sprintf("Could not parse value: '%s', integer expected", s)
	}
	return value, msg
}

func (t *ParamTransformer) TransformSplit(s string) ([]string, string) {
	return strings.Split(s, ","), ""
}

func (t *ParamTransformer) TransformSplitAndUppercase(s string) ([]string, string) {
	sList, _ := t.TransformSplit(s)
	for i := 0; i < len(sList); i++ {
		sList[i] = strings.ToUpper(sList[i])
	}
	return sList, ""
}
