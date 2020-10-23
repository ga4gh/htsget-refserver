// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module set_test tests module set
package htsrequest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// setSingleParameterTC test cases for SetSingleParameter
var setSingleParamterTC = []struct {
	endpoint        htsconstants.APIEndpoint
	method          string
	id              string
	queryString     string
	headers         [][]string
	requestBody     string
	setParamTuple   *SetParameterTuple
	expError        bool
	expErrorMessage string
}{
	{
		// error when validating id
		htsconstants.APIEndpointReadsTicket,
		"GET",
		"NoID",
		"",
		[][]string{},
		"",
		&SetParameterTuple{
			htsconstants.ParamLocPath,
			"id",
			"NoTransform",
			"ValidateID",
			"SetID",
			defaultID,
		},
		true,
		"The requested resource could not be associated with a registered data source",
	},
	{
		// successful parse of format
		htsconstants.APIEndpointReadsTicket,
		"GET",
		"NoID",
		"?format=BAM",
		[][]string{},
		"",
		&SetParameterTuple{
			htsconstants.ParamLocQuery,
			"format",
			"TransformStringUppercase",
			"ValidateFormat",
			"SetFormat",
			defaultFormatReads,
		},
		false,
		"",
	},
	{
		// error when validating format
		htsconstants.APIEndpointReadsTicket,
		"GET",
		"NoID",
		"?format=VCF",
		[][]string{},
		"",
		&SetParameterTuple{
			htsconstants.ParamLocQuery,
			"format",
			"TransformStringUppercase",
			"ValidateFormat",
			"SetFormat",
			defaultFormatReads,
		},
		true,
		"file format: 'VCF' not supported",
	},
	{
		// use default format
		htsconstants.APIEndpointReadsTicket,
		"GET",
		"NoID",
		"",
		[][]string{},
		"",
		&SetParameterTuple{
			htsconstants.ParamLocQuery,
			"format",
			"TransformStringUppercase",
			"ValidateFormat",
			"SetFormat",
			defaultFormatReads,
		},
		false,
		"",
	},
	{
		// successful parse of HtsgetTotalBlocks
		htsconstants.APIEndpointReadsTicket,
		"GET",
		"NoID",
		"",
		[][]string{[]string{"HtsgetTotalBlocks", "100"}},
		"",
		&SetParameterTuple{
			htsconstants.ParamLocHeader,
			"HtsgetTotalBlocks",
			"NoTransform",
			"NoValidation",
			"SetHtsgetTotalBlocks",
			defaultHtsgetTotalBlocks,
		},
		false,
		"",
	},
	{
		// parse fields during POST request
		htsconstants.APIEndpointReadsTicket,
		"POST",
		"NoID",
		"",
		[][]string{},
		"{\"fields\": [\"QNAME\",\"SEQ\",\"QUAL\"]}",
		&SetParameterTuple{
			htsconstants.ParamLocReqBody,
			"fields",
			"NoTransform",
			"ValidateFields",
			"SetFields",
			defaultFields,
		},
		false,
		"",
	},
}

// TestSetSingleParameter tests SetSingleParameter function
func TestSetSingleParameter(t *testing.T) {
	for _, tc := range setSingleParamterTC {

		// setup requestBody
		var requestBodyBytes []byte = []byte{}
		if tc.requestBody != "" {
			requestBodyBytes = []byte(tc.requestBody)
		}

		// router handler function runs SetSingleParameter based on test case,
		// and asserts whether an expected error is present or not
		routerHandler := func(writer http.ResponseWriter, request *http.Request) {
			// setup base htsgetRequest
			htsgetReq := NewHtsgetRequest()
			htsgetReq.SetEndpoint(tc.endpoint)

			// run setSingleParameter function, validate error either nil or not nil
			// as expected
			err := setSingleParameter(request, *tc.setParamTuple, requestBodyBytes, htsgetReq)
			if tc.expError {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expErrorMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
		}

		// configure route handler on both GET and POST
		router := chi.NewRouter()
		router.Get("/test/{id}", routerHandler)
		router.Post("/test/{id}", routerHandler)

		// setup HTTP request, add headers to request if applicable
		request := httptest.NewRequest(tc.method, "/test/"+tc.id+tc.queryString, nil)
		if len(tc.headers) > 0 {
			for _, header := range tc.headers {
				request.Header.Add(header[0], header[1])
			}
		}

		// execute the single request, which will trigger the above handler function
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, request)
	}
}
