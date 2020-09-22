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

type RequestParameter struct {
	location htsconstants.ParamLoc
	name     string
}

var paramsByMethodEndpointLocation = map[htsconstants.HTTPMethod]map[htsconstants.APIEndpoint][]RequestParameter{
	htsconstants.GetMethod: map[htsconstants.APIEndpoint][]RequestParameter{
		htsconstants.APIEndpointReadsTicket: []RequestParameter{
			{htsconstants.ParamLocPath, "id"},
			{htsconstants.ParamLocQuery, "format"},
			{htsconstants.ParamLocQuery, "class"},
			{htsconstants.ParamLocQuery, "referenceName"},
			{htsconstants.ParamLocQuery, "start"},
			{htsconstants.ParamLocQuery, "end"},
			{htsconstants.ParamLocQuery, "fields"},
			{htsconstants.ParamLocQuery, "tags"},
			{htsconstants.ParamLocQuery, "notags"},
		},
		htsconstants.APIEndpointReadsData: []RequestParameter{
			{htsconstants.ParamLocPath, "id"},
			{htsconstants.ParamLocQuery, "format"},
			{htsconstants.ParamLocQuery, "referenceName"},
			{htsconstants.ParamLocQuery, "start"},
			{htsconstants.ParamLocQuery, "end"},
			{htsconstants.ParamLocQuery, "fields"},
			{htsconstants.ParamLocQuery, "tags"},
			{htsconstants.ParamLocQuery, "notags"},
			{htsconstants.ParamLocHeader, "HtsgetBlockClass"},
			{htsconstants.ParamLocHeader, "HtsgetBlockId"},
			{htsconstants.ParamLocHeader, "HtsgetNumBlocks"},
		},
		htsconstants.APIEndpointReadsServiceInfo: []RequestParameter{},
		htsconstants.APIEndpointVariantsTicket: []RequestParameter{
			{htsconstants.ParamLocPath, "id"},
			{htsconstants.ParamLocQuery, "format"},
			{htsconstants.ParamLocQuery, "class"},
			{htsconstants.ParamLocQuery, "referenceName"},
			{htsconstants.ParamLocQuery, "start"},
			{htsconstants.ParamLocQuery, "end"},
			{htsconstants.ParamLocQuery, "fields"},
			{htsconstants.ParamLocQuery, "tags"},
			{htsconstants.ParamLocQuery, "notags"},
		},
		htsconstants.APIEndpointVariantsData: []RequestParameter{
			{htsconstants.ParamLocPath, "id"},
			{htsconstants.ParamLocQuery, "format"},
			{htsconstants.ParamLocQuery, "class"},
			{htsconstants.ParamLocQuery, "referenceName"},
			{htsconstants.ParamLocQuery, "start"},
			{htsconstants.ParamLocQuery, "end"},
			{htsconstants.ParamLocQuery, "fields"},
			{htsconstants.ParamLocQuery, "tags"},
			{htsconstants.ParamLocQuery, "notags"},
			{htsconstants.ParamLocHeader, "HtsgetBlockClass"},
			{htsconstants.ParamLocHeader, "HtsgetBlockId"},
			{htsconstants.ParamLocHeader, "HtsgetNumBlocks"},
		},
		htsconstants.APIEndpointFileBytes: []RequestParameter{
			{htsconstants.ParamLocHeader, "HtsgetFilePath"},
			{htsconstants.ParamLocHeader, "Range"},
		},
	},
	htsconstants.PostMethod: map[htsconstants.APIEndpoint][]RequestParameter{
		htsconstants.APIEndpointReadsTicket: []RequestParameter{
			{htsconstants.ParamLocPath, "id"},
			{htsconstants.ParamLocReqBody, "referenceName"},
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
func setSingleParameter(request *http.Request, param RequestParameter,
	params url.Values, htsgetReq *HtsgetRequest) error {

	var value string
	var found bool
	// lookup if parameter is found on path/query/header,
	// and if a scalar or list is expected
	paramLocation := param.location
	paramName := param.name
	paramType := paramTypes[paramName]

	// parse the request parameter by path, query string, or header
	switch paramLocation {
	case htsconstants.ParamLocPath:
		value, found = parsePathParam(request, paramName)
	case htsconstants.ParamLocQuery:
		v, f, err := parseQueryParam(params, paramName)
		value = v
		found = f
		if err != nil {
			return err
		}
	case htsconstants.ParamLocHeader:
		value, found = parseHeaderParam(request, paramName)
	case htsconstants.ParamLocReqBody:
		value, found = parseReqBodyParam(request, paramName)
	}

	// if a value is found, then
	if found {
		// run the validation function, return an error if invalid
		validationFunc := validationByParam[paramName]
		validationResult, validationMsg := validationFunc(value, htsgetReq)
		if !validationResult {
			return errors.New(validationMsg)
		}

		// if valid, transform the param value and set it to the
		// HtsgetRequest map
		switch paramType {
		case ParamTypeScalar:
			transformFunc := transformationScalarByParam[paramName]
			htsgetReq.AddScalarParam(paramName, transformFunc(value))
		case ParamTypeList:
			transformFunc := transformationListByParam[paramName]
			htsgetReq.AddListParam(paramName, transformFunc(value))
		}
		return nil
	}

	// if no param value is found, set the default value to the HtsgetRequest
	// map
	switch paramType {
	case ParamTypeScalar:
		htsgetReq.AddScalarParam(paramName, defaultScalarParameterValues[paramName])
	case ParamTypeList:
		htsgetReq.AddListParam(paramName, defaultListParameterValues[paramName])
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

	orderedParams := paramsByMethodEndpointLocation[method][endpoint]
	htsgetReq := NewHtsgetRequest()
	htsgetReq.SetEndpoint(endpoint)
	params := request.URL.Query()
	for i := 0; i < len(orderedParams); i++ {
		param := orderedParams[i]
		paramName := param.name
		err := setSingleParameter(request, param, params, htsgetReq)
		if err != nil {
			htsgetErrorFunc := errorsByParam[paramName]
			msg := err.Error()
			htsgetErrorFunc(writer, &msg)
			return htsgetReq, err
		}
	}
	return htsgetReq, nil
}
