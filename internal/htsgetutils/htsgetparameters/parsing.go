package htsgetparameters

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

var isPathByParam = map[string]bool{
	"id":            true,
	"format":        false,
	"class":         false,
	"referenceName": false,
	"start":         false,
	"end":           false,
	"fields":        false,
	"tags":          false,
	"notags":        false,
}

var isScalarByParam = map[string]bool{
	"id":            true,
	"format":        true,
	"class":         true,
	"referenceName": true,
	"start":         true,
	"end":           true,
	"fields":        false,
	"tags":          false,
	"notags":        false,
}

func parsePathParam(request *http.Request, key string) string {
	value := chi.URLParam(request, key)
	return value
}

func parseQueryParam(params url.Values, key string) (string, error) {
	if len(params[key]) == 1 {
		return params[key][0], nil
	}
	if len(params[key]) > 1 {
		return "", errors.New("too many values specified for parameter: " + key)
	}
	return "", nil
}
