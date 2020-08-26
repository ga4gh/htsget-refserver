// Package htserror write HTTP errors with various codes, response bodies
//
// Modules error_test tests error
package htserror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var newMessageA = "format BAM not supported"

var errorsTC = []struct {
	errorFunction func(writer http.ResponseWriter, msgPtr *string)
	message       *string
	expString     string
	expCode       int
}{
	{
		UnsupportedFormat,
		&newMessageA,
		"UnsupportedFormat: format BAM not supported",
		codeBadRequest,
	},
	{
		UnsupportedFormat,
		nil,
		"UnsupportedFormat: The requested file format is not supported by the server",
		codeBadRequest,
	},
	{
		InvalidInput,
		nil,
		"InvalidInput: The request parameters do not adhere to the specification",
		codeBadRequest,
	},
	{
		InvalidRange,
		nil,
		"InvalidRange: The requested range cannot be satisfied",
		codeBadRequest,
	},
	{
		InvalidAuthentication,
		nil,
		"InvalidAuthentication: Authorization provided is invalid",
		codeInvalidAuthentication,
	},
	{
		PermissionDenied,
		nil,
		"PermissionDenied: Authorization is required to access the resource",
		codePermissionDenied,
	},
	{
		NotFound,
		nil,
		"NotFound: The resource requested was not found",
		codeNotFound,
	},
	{
		InternalServerError,
		nil,
		"InternalServerError: Internal server error",
		codeInternalServerError,
	},
}

func TestErrors(t *testing.T) {
	for _, tc := range errorsTC {

		writer := httptest.NewRecorder()
		tc.errorFunction(writer, tc.message)

		bytes := make([]byte, 512)
		nBytes, _ := writer.Body.Read(bytes)
		bytes = bytes[0:nBytes]
		htsgetErrObj := new(htsgetError)
		json.Unmarshal(bytes, htsgetErrObj)

		assert.Equal(t, tc.expString, htsgetErrObj.Error())
		assert.Equal(t, tc.expCode, writer.Code)
	}
}
