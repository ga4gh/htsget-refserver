// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module set defines operations for setting request parameters to an
// HtsgetRequest, which first involves correct parsing, validation, and
// transformation. Sets parameters correctly based on request route
package htsrequest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/ga4gh/htsget-refserver/internal/htserror"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// SetParameterTuple describes how a single request parameter will be parsed,
// transformed, validated, and then set to the mature htsgetrequest object
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
			{
				htsconstants.ParamLocReqBody,
				"format",
				"NoTransform",
				"ValidateFormat",
				"SetFormat",
			},
			{
				htsconstants.ParamLocReqBody,
				"fields",
				"NoTransform",
				"ValidateFields",
				"SetFields",
			},
			{
				htsconstants.ParamLocReqBody,
				"tags",
				"NoTransform",
				"ValidateTags",
				"SetTags",
			},
			{
				htsconstants.ParamLocReqBody,
				"notags",
				"NoTransform",
				"ValidateNoTags",
				"SetNoTags",
			},
			{
				htsconstants.ParamLocReqBody,
				"regions",
				"NoTransform",
				"ValidateRegions",
				"SetRegions",
			},
		},
	},
}

// setSingleParameter parses, transforms, validates, and sets a valid parameter
// to the HtsgetRequest object. if the parameter value is not valid,
// returns an error
func setSingleParameter(request *http.Request, setParamTuple SetParameterTuple,
	requestBodyBytes []byte, htsgetReq *HtsgetRequest) error {

	var rawValue string              // raw string value parsed from query string, path, or header
	var reflectedValue reflect.Value // post-transform, post-reflection representation of param
	var found bool

	// lookup if parameter is found on path/query/header,
	// and if a scalar or list is expected
	location := setParamTuple.location
	paramName := setParamTuple.name

	// parse the request parameter by path, query string, header, or request body
	switch location {
	case htsconstants.ParamLocPath:
		rawValue, found = parsePathParam(request, paramName)
	case htsconstants.ParamLocQuery:
		v, f, err := parseQueryParam(request.URL.Query(), paramName)
		rawValue = v
		found = f
		if err != nil {
			return err
		}
	case htsconstants.ParamLocHeader:
		rawValue, found = parseHeaderParam(request, paramName)
	case htsconstants.ParamLocReqBody:
		v, f, err := parseReqBodyParam(requestBodyBytes, paramName)
		reflectedValue = v
		found = f
		if err != nil {
			return err
		}
	}

	// use reflect to get the param setter method for the request
	htsgetReqReflect := reflect.ValueOf(htsgetReq)
	htsgetParamSetter := htsgetReqReflect.MethodByName(setParamTuple.setFunc)

	// if the value was found in path, query, header, or req body
	if found {

		// path, query, and header params must be transformed from a string to another
		// datatype (if necessary), and then reflected via reflect API
		// req body params do not undergo this step as they are inherently in their
		// own datatype, and have been reflected
		if location == htsconstants.ParamLocPath || location == htsconstants.ParamLocQuery || location == htsconstants.ParamLocHeader {
			// use reflection to call the transformation function by name
			transformer := NewParamTransformer()
			transformerReflect := reflect.ValueOf(transformer)
			transformFunc := transformerReflect.MethodByName(setParamTuple.transformFunc)
			transformResult := transformFunc.Call([]reflect.Value{reflect.ValueOf(rawValue)})
			reflectedValue = transformResult[0]
			message := transformResult[1].String()
			if message != "" {
				return errors.New(message)
			}
		}

		// use reflection to call the validation function by name
		validator := NewParamValidator()
		validatorReflect := reflect.ValueOf(validator)
		validateFunc := validatorReflect.MethodByName(setParamTuple.validateFunc)
		resultMsg := validateFunc.Call([]reflect.Value{reflect.ValueOf(htsgetReq), reflectedValue})
		result := resultMsg[0].Bool()
		msg := resultMsg[1].String()
		if !result {
			return errors.New(msg)
		}

		// if validation passed, set the transformed value
		htsgetParamSetter.Call([]reflect.Value{reflectedValue})
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

	// for POST requests, unmarshal the JSON body once and pass to individual
	// setting methods
	var requestBodyBytes []byte
	if method == htsconstants.PostMethod {
		rbb, err := ioutil.ReadAll(request.Body)
		requestBodyBytes = rbb
		msg := "Request body malformed"
		if err != nil {
			htserror.InvalidInput(writer, &msg)
			return htsgetReq, err
		}
	}

	for i := 0; i < len(orderedParams); i++ {
		param := orderedParams[i]
		paramName := param.name
		err := setSingleParameter(request, param, requestBodyBytes, htsgetReq)
		if err != nil {
			htsgetErrorFunc := errorsByParam[paramName]
			msg := err.Error()
			htsgetErrorFunc(writer, &msg)
			return htsgetReq, err
		}
	}
	return htsgetReq, nil
}
