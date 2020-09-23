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
	"reflect"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

type SetParameterTuple struct {
	location      htsconstants.ParamLoc
	name          string
	transformFunc string
	validateFunc  string
	setFunc       string
}

var orderedParamsMap = map[htsconstants.HTTPMethod]map[htsconstants.APIEndpoint][]SetParameterTuple{

	/* **************************************************
	 * HTTP GET
	 * ************************************************** */

	htsconstants.GetMethod: map[htsconstants.APIEndpoint][]SetParameterTuple{

		/* **************************************************
		 * HTTP GET READS TICKET
		 * ************************************************** */

		htsconstants.APIEndpointReadsTicket: []SetParameterTuple{
			{
				htsconstants.ParamLocPath,
				"id",
				"NoTransform",
				"ValidateID",
				"SetID",
			},
			{
				htsconstants.ParamLocQuery,
				"format",
				"TransformStringUppercase",
				"ValidateFormat",
				"SetFormat",
			},

			{
				htsconstants.ParamLocQuery,
				"class",
				"TransformStringLowercase",
				"ValidateClass",
				"SetClass",
			},
			{
				htsconstants.ParamLocQuery,
				"referenceName",
				"NoTransform",
				"ValidateReferenceName",
				"SetReferenceName",
			},
			{
				htsconstants.ParamLocQuery,
				"start",
				"TransformStringToInt",
				"ValidateStart",
				"SetStart",
			},
			{
				htsconstants.ParamLocQuery,
				"end",
				"TransformStringToInt",
				"ValidateEnd",
				"SetEnd",
			},
			{
				htsconstants.ParamLocQuery,
				"fields",
				"TransformSplitAndUppercase",
				"ValidateFields",
				"SetFields",
			},
			{
				htsconstants.ParamLocQuery,
				"tags",
				"TransformSplit",
				"ValidateTags",
				"SetTags",
			},
			{
				htsconstants.ParamLocQuery,
				"notags",
				"TransformSplit",
				"ValidateNoTags",
				"SetNoTags",
			},
		},

		/* **************************************************
		 * HTTP GET READS DATA
		 * ************************************************** */

		htsconstants.APIEndpointReadsData: []SetParameterTuple{
			{
				htsconstants.ParamLocPath,
				"id",
				"NoTransform",
				"ValidateID",
				"SetID",
			},
			{
				htsconstants.ParamLocQuery,
				"format",
				"TransformStringUppercase",
				"ValidateFormat",
				"SetFormat",
			},
			{
				htsconstants.ParamLocQuery,
				"referenceName",
				"NoTransform",
				"ValidateReferenceName",
				"SetReferenceName",
			},
			{
				htsconstants.ParamLocQuery,
				"start",
				"TransformStringToInt",
				"ValidateStart",
				"SetStart",
			},
			{
				htsconstants.ParamLocQuery,
				"end",
				"TransformStringToInt",
				"ValidateEnd",
				"SetEnd",
			},
			{
				htsconstants.ParamLocQuery,
				"fields",
				"TransformSplitAndUppercase",
				"ValidateFields",
				"SetFields",
			},
			{
				htsconstants.ParamLocQuery,
				"tags",
				"TransformSplit",
				"ValidateTags",
				"SetTags",
			},
			{
				htsconstants.ParamLocQuery,
				"notags",
				"TransformSplit",
				"ValidateNoTags",
				"SetNoTags",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetBlockClass",
				"TransformStringLowercase",
				"NoValidation",
				"SetHtsgetBlockClass",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetCurrentBlock",
				"NoTransform",
				"NoValidation",
				"SetHtsgetCurrentBlock",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetTotalBlocks",
				"NoTransform",
				"NoValidation",
				"SetHtsgetTotalBlocks",
			},
		},

		/* **************************************************
		 * HTTP GET READS SERVICE INFO
		 * ************************************************** */

		htsconstants.APIEndpointReadsServiceInfo: []SetParameterTuple{},

		/* **************************************************
		 * HTTP GET VARIANTS TICKET
		 * ************************************************** */

		htsconstants.APIEndpointVariantsTicket: []SetParameterTuple{
			{
				htsconstants.ParamLocPath,
				"id",
				"NoTransform",
				"ValidateID",
				"SetID",
			},
			{
				htsconstants.ParamLocQuery,
				"format",
				"TransformStringUppercase",
				"ValidateFormat",
				"SetFormat",
			},

			{
				htsconstants.ParamLocQuery,
				"class",
				"TransformStringLowercase",
				"ValidateClass",
				"SetClass",
			},
			{
				htsconstants.ParamLocQuery,
				"referenceName",
				"NoTransform",
				"ValidateReferenceName",
				"SetReferenceName",
			},
			{
				htsconstants.ParamLocQuery,
				"start",
				"TransformStringToInt",
				"ValidateStart",
				"SetStart",
			},
			{
				htsconstants.ParamLocQuery,
				"end",
				"TransformStringToInt",
				"ValidateEnd",
				"SetEnd",
			},
			{
				htsconstants.ParamLocQuery,
				"fields",
				"TransformSplitAndUppercase",
				"ValidateFields",
				"SetFields",
			},
			{
				htsconstants.ParamLocQuery,
				"tags",
				"TransformSplit",
				"ValidateTags",
				"SetTags",
			},
			{
				htsconstants.ParamLocQuery,
				"notags",
				"TransformSplit",
				"ValidateNoTags",
				"SetNoTags",
			},
		},

		/* **************************************************
		 * HTTP GET VARIANTS DATA
		 * ************************************************** */

		htsconstants.APIEndpointVariantsData: []SetParameterTuple{
			{
				htsconstants.ParamLocPath,
				"id",
				"NoTransform",
				"ValidateID",
				"SetID",
			},
			{
				htsconstants.ParamLocQuery,
				"format",
				"TransformStringUppercase",
				"ValidateFormat",
				"SetFormat",
			},
			{
				htsconstants.ParamLocQuery,
				"referenceName",
				"NoTransform",
				"ValidateReferenceName",
				"SetReferenceName",
			},
			{
				htsconstants.ParamLocQuery,
				"start",
				"TransformStringToInt",
				"ValidateStart",
				"SetStart",
			},
			{
				htsconstants.ParamLocQuery,
				"end",
				"TransformStringToInt",
				"ValidateEnd",
				"SetEnd",
			},
			{
				htsconstants.ParamLocQuery,
				"fields",
				"TransformSplitAndUppercase",
				"ValidateFields",
				"SetFields",
			},
			{
				htsconstants.ParamLocQuery,
				"tags",
				"TransformSplit",
				"ValidateTags",
				"SetTags",
			},
			{
				htsconstants.ParamLocQuery,
				"notags",
				"TransformSplit",
				"ValidateNoTags",
				"SetNoTags",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetBlockClass",
				"TransformStringLowercase",
				"NoValidation",
				"SetHtsgetBlockClass",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetCurrentBlock",
				"NoTransform",
				"NoValidation",
				"SetHtsgetCurrentBlock",
			},
			{
				htsconstants.ParamLocHeader,
				"HtsgetTotalBlocks",
				"NoTransform",
				"NoValidation",
				"SetHtsgetTotalBlocks",
			},
		},

		/* **************************************************
		 * HTTP GET VARIANTS SERVICE INFO
		 * ************************************************** */

		htsconstants.APIEndpointVariantsServiceInfo: []SetParameterTuple{},

		/* **************************************************
		 * HTTP GET FILE BYTES
		 * ************************************************** */

		htsconstants.APIEndpointFileBytes: []SetParameterTuple{
			{
				htsconstants.ParamLocHeader,
				"HtsgetFilePath",
				"NoTransform",
				"NoValidation",
				"SetHtsgetFilePath",
			},
			{
				htsconstants.ParamLocHeader,
				"Range",
				"NoTransform",
				"NoValidation",
				"SetHtsgetRange",
			},
		},
	},

	/* **************************************************
	 * HTTP POST
	 * ************************************************** */

	htsconstants.PostMethod: map[htsconstants.APIEndpoint][]SetParameterTuple{

		/* **************************************************
		 * HTTP POST READS TICKET
		 * ************************************************** */

		htsconstants.APIEndpointReadsTicket: []SetParameterTuple{
			{
				htsconstants.ParamLocPath,
				"id",
				"NoTransform",
				"ValidateID",
				"SetID",
			},
		},
	},
}

// setSingleParameter parses, transforms, validates, and sets a valid parameter
// to the HtsgetRequest object. if the parameter value is not valid,
// returns an error
func setSingleParameter(request *http.Request, setParamTuple SetParameterTuple,
	params url.Values, htsgetReq *HtsgetRequest) error {

	var value string
	var found bool
	// lookup if parameter is found on path/query/header,
	// and if a scalar or list is expected
	location := setParamTuple.location
	paramName := setParamTuple.name

	// parse the request parameter by path, query string, or header
	switch location {
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

	// use reflect to get the param setter method for the request
	htsgetReqReflect := reflect.ValueOf(htsgetReq)
	htsgetParamSetter := htsgetReqReflect.MethodByName(setParamTuple.setFunc)

	// if a value is found, then transform, validate, and set
	if found {
		// use reflection to call the transformation function by name
		transformer := NewParamTransformer()
		transformerReflect := reflect.ValueOf(transformer)
		transformFunc := transformerReflect.MethodByName(setParamTuple.transformFunc)
		transformResult := transformFunc.Call([]reflect.Value{reflect.ValueOf(value)})
		transformed := transformResult[0]
		message := transformResult[1].String()
		if message != "" {
			return errors.New(message)
		}

		// use reflection to call the validation function by name
		validator := NewParamValidator()
		validatorReflect := reflect.ValueOf(validator)
		validateFunc := validatorReflect.MethodByName(setParamTuple.validateFunc)
		resultMsg := validateFunc.Call([]reflect.Value{reflect.ValueOf(htsgetReq), transformed})
		result := resultMsg[0].Bool()
		message = resultMsg[1].String()
		if !result {
			return errors.New(message)
		}

		// if validation passed, set the transformed value
		htsgetParamSetter.Call([]reflect.Value{transformed})
		return nil
	}

	// if no param value is found, set the default value
	defaultValueReflect := reflect.ValueOf(defaultParameterValues[paramName])
	htsgetParamSetter.Call([]reflect.Value{defaultValueReflect})
	return nil
}

// SetAllParameters parses, transforms, validates, and sets all parameters to
// an HtsgetRequest for a given ordered list of expected request parameters
func SetAllParameters(method htsconstants.HTTPMethod, endpoint htsconstants.APIEndpoint, writer http.ResponseWriter, request *http.Request) (*HtsgetRequest, error) {

	orderedParams := orderedParamsMap[method][endpoint]
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
