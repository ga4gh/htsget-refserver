// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module set.go defines operations for setting request parameters to an
// HtsgetRequest, which first involves correct parsing, validation, and
// transformation. Sets parameters correctly based on request route
package htsrequest

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

var orderedParametersByMethodAndEndpoint = map[htsconstants.HTTPMethod]map[htsconstants.APIEndpoint][]string{
	htsconstants.GetMethod: map[htsconstants.APIEndpoint][]string{
		htsconstants.APIEndpointReadsTicket: []string{
			"id",
			"format",
			"class",
			"referenceName",
			"start",
			"end",
			"fields",
			"tags",
			"notags",
		},
		htsconstants.APIEndpointReadsData: []string{
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
		},
		htsconstants.APIEndpointReadsServiceInfo: []string{},
		htsconstants.APIEndpointVariantsTicket: []string{
			"id",
			"format",
			"class",
			"referenceName",
			"start",
			"end",
			"fields",
			"tags",
			"notags",
		},
		htsconstants.APIEndpointVariantsData: []string{
			"id",
			"format",
			"class",
			"referenceName",
			"start",
			"end",
			"fields",
			"tags",
			"notags",
			"HtsgetBlockClass",
			"HtsgetBlockId",
			"HtsgetNumBlocks",
		},
		htsconstants.APIEndpointFileBytes: []string{
			"HtsgetFilePath",
			"Range",
		},
	},
}

// setSingleParameter parses, validates, and sets a valid parameter to the
// HtsgetRequest object. if the parameter value is not valid, returns an error
//
// Arguments
//	request (*http.Request): HTTP request object
//	paramKey (string): parameter name to parse, validate, etc.
//	params (url.Values): query string parameters from HTTP request
// 	htsgetReq (*HtsgetRequest): object to set transformed parameter value to
// Returns
//	(error): client-side error if any parameters fail validation
func setSingleParameter(request *http.Request, paramKey string,
	params url.Values, htsgetReq *HtsgetRequest) error {

	var value string
	var found bool
	// lookup if parameter is found on path/query/header,
	// and if a scalar or list is expected
	paramLocation := paramLocations[paramKey]
	paramType := paramTypes[paramKey]

	// parse the request parameter by path, query string, or header
	switch paramLocation {
	case ParamLocPath:
		value, found = parsePathParam(request, paramKey)
	case ParamLocQuery:
		v, f, err := parseQueryParam(params, paramKey)
		value = v
		found = f
		if err != nil {
			return err
		}
	case ParamLocHeader:
		value, found = parseHeaderParam(request, paramKey)
	}

	// if a value is found, then
	if found {
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

// SetAllParameters parses, validates, transforms, and sets all parameters to
// an HtsgetRequest for a given ordered list of expected request parameters
//
// Arguments
//	orderedParams ([]string): route-specific order to parse parameters in
//	request (*http.Request): HTTP request
//	writer (http.ResponseWriter): HTTP response writer (to write error if necessary)
//	params (url.Values): query string parameters
// Returns
//	(*HtsgetRequest): object with mature parameters set to it
//	(error): client-side error if any parameters fail validation
func SetAllParameters(method htsconstants.HTTPMethod, endpoint htsconstants.APIEndpoint, writer http.ResponseWriter, request *http.Request) (*HtsgetRequest, error) {

	orderedParams := orderedParametersByMethodAndEndpoint[method][endpoint]
	htsgetReq := NewHtsgetRequest()
	htsgetReq.SetEndpoint(endpoint)
	params := request.URL.Query()
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
