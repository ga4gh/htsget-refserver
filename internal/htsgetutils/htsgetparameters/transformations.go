package htsgetparameters

import (
	"strings"
)

var transformationScalarByParam = map[string]func(string) string{
	"id":            noTransform,
	"format":        strings.ToUpper,
	"class":         strings.ToLower,
	"referenceName": noTransform,
	"start":         noTransform,
	"end":           noTransform,
}

var transformationListByParam = map[string]func(string) []string{
	"fields": splitAndUppercase,
	"tags":   splitOnComma,
	"notags": splitOnComma,
}

func noTransform(s string) string {
	return s
}

func splitOnComma(s string) []string {
	return strings.Split(s, ",")
}

func splitAndUppercase(s string) []string {
	sList := splitOnComma(s)
	for i := 0; i < len(sList); i++ {
		sList[i] = strings.ToUpper(sList[i])
	}
	return sList
}
