// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

import (
	"errors"
	"net/http"
	"net/url"
)

// the order in which parameters (path or query) are parsed, validated, and
// set in. the validation of parameter later in the order may be dependent
// on previously validated parameters (e.g. 'tags' and 'notags')
var readsTicketEndpointSetParamsOrder = []string{
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

// for the  the order in which parameters (path or query) are parsed, validated, and set in
var readsDataEndpointSetParamsOrder = []string{
	"id",
	"format",
	"referenceName",
	"start",
	"end",
	"fields",
	"tags",
	"notags",
	"HtsgetBlockClass",
	"HtsgetBlockId",
	"HtsgetNumBlocks",
}

// parses, validates, and sets valid parameter to the HtsgetRequest object.
// if the parameter value is not valid, returns an error
func setSingleParameter(request *http.Request, paramKey string,
	params url.Values, htsgetReq *HtsgetRequest) error {

	var value string
	// map lookup to determine if parameter is found on path/query string,
	// and if a scalar or list is expected
	paramLocation := paramLocations[paramKey]
	paramType := paramTypes[paramKey]

	// parse the request parameter by path, query string, or header
	switch paramLocation {
	case ParamLocPath:
		value = parsePathParam(request, paramKey)
	case ParamLocQuery:
		v, err := parseQueryParam(params, paramKey)
		value = v
		if err != nil {
			return err
		}
	case ParamLocHeader:
		value = parseHeaderParam(request, paramKey)
	}

	// if a value is found, then
	if value != "" {
		// run the validation function, return an error if invalid
		validationFunc := validationByParam[paramKey]
		validationResult, validationMsg := validationFunc(value, htsgetReq)
		if !validationResult {
			return errors.New(validationMsg)
		}

		// if valid, transform the param value and set it to the
		// HtsgetRequest map
		switch paramType {
		case ParamTypeScalar:
			transformFunc := transformationScalarByParam[paramKey]
			htsgetReq.AddScalarParam(paramKey, transformFunc(value))
		case ParamTypeList:
			transformFunc := transformationListByParam[paramKey]
			htsgetReq.AddListParam(paramKey, transformFunc(value))
		}
		return nil
	}

	// if no param value is found, set the default value to the HtsgetRequest
	// map
	switch paramType {
	case ParamTypeScalar:
		htsgetReq.AddScalarParam(paramKey, defaultScalarParameterValues[paramKey])
	case ParamTypeList:
		htsgetReq.AddListParam(paramKey, defaultListParameterValues[paramKey])
	}
	return nil
}

// setAllParameters for a given ordered list of expected request parameters,
// parse, set, validate, and transform all parameters in order, setting them
// to an HtsgetRequest
func setAllParameters(orderedParams []string, request *http.Request, writer http.ResponseWriter, params url.Values) (*HtsgetRequest, error) {
	htsgetReq := NewHtsgetRequest()

	for i := 0; i < len(orderedParams); i++ {
		paramKey := orderedParams[i]
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

// ReadsTicketEndpointSetAllParameters sets all parameters expected in the
// reads ticket endpoint to an HtsgetRequest
func ReadsTicketEndpointSetAllParameters(request *http.Request, writer http.ResponseWriter, params url.Values) (*HtsgetRequest, error) {
	return setAllParameters(readsTicketEndpointSetParamsOrder, request, writer, params)
}

// ReadsDataEndpointSetAllParameters sets all parameters expected in the
// reads data endpoint to an HtsgetRequest
func ReadsDataEndpointSetAllParameters(request *http.Request, writer http.ResponseWriter, params url.Values) (*HtsgetRequest, error) {
	return setAllParameters(readsDataEndpointSetParamsOrder, request, writer, params)
}
