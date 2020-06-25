package htsgetparameters

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ga4gh/htsget-refserver/internal/htsgetutils/htsgethttp/htsgetrequest"
)

var readsEndpointSetParamsOrder = []string{
	"id",
	"format",
	"class",
	"referenceName",
	"start",
	"end",
	"fields",
	"tags",
	"notags",
}

func setSingleParameter(request *http.Request, paramKey string,
	params url.Values, htsgetReq *htsgetrequest.HtsgetRequest) error {

	var value string
	isPath := isPathByParam[paramKey]
	isScalar := isScalarByParam[paramKey]

	if isPath {
		value = parsePathParam(request, paramKey)
	} else {
		v, err := parseQueryParam(params, paramKey)
		value = v
		if err != nil {
			return err
		}
	}

	if value != "" {
		validationFunc := validationByParam[paramKey]
		validationResult, validationMsg := validationFunc(value, htsgetReq)
		if !validationResult {
			return errors.New(validationMsg)
		}

		if isScalar {
			transformFunc := transformationScalarByParam[paramKey]
			htsgetReq.AddToScalars(paramKey, transformFunc(value))
		} else {
			transformFunc := transformationListByParam[paramKey]
			htsgetReq.AddToLists(paramKey, transformFunc(value))
		}

		return nil
	}

	if isScalar {
		htsgetReq.AddToScalars(paramKey, defaultScalarParameterValues[paramKey])
	} else {
		htsgetReq.AddToLists(paramKey, defaultListParameterValues[paramKey])
	}

	return nil
}

func ReadsEndpointSetAllParameters(request *http.Request, writer http.ResponseWriter, params url.Values) (*htsgetrequest.HtsgetRequest, error) {
	htsgetReq := htsgetrequest.New()

	for i := 0; i < len(readsEndpointSetParamsOrder); i++ {
		paramKey := readsEndpointSetParamsOrder[i]
		err := setSingleParameter(request, paramKey, params, htsgetReq)
		if err != nil {
			htsgetErrorFunc := errorsByParam[paramKey]
			msg := err.Error()
			htsgetErrorFunc(writer, &msg)
			return htsgetReq, err
		}
	}
	return htsgetReq, nil
}
