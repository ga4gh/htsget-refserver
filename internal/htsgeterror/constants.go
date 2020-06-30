// Package htsgeterror provides operations for writing various HTTP errors,
// including different status codes and error response bodies
//
// Module constants.go contains HTTP Error constants, including status codes,
// error names, and default messages
package htsgeterror

import (
	"net/http"
	"strconv"
)

/* Error Status Codes: HTTP status code constants by error */

// codeBadRequest (int): status code for Bad Request errors
const codeBadRequest = http.StatusBadRequest

// codeInvalidAuthentication (int): status code for AuthN errors
const codeInvalidAuthentication = http.StatusUnauthorized

// codePermissionDenied (int): status code for AuthZ errors
const codePermissionDenied = http.StatusForbidden

// codeNotFound (int): status code for Not Found errors
const codeNotFound = http.StatusNotFound

// codeInternalServerError (int): status code for unspecified server-side error
const codeInternalServerError = http.StatusInternalServerError

/* Error Names: htsget canonical error names */

// errorBadRequestUnsupportedFormat (string): error name for unsupported format
const errorBadRequestUnsupportedFormat = "UnsupportedFormat"

// errorBadRequestInvalidInput (string): error name for invalid input
const errorBadRequestInvalidInput = "InvalidInput"

// errorBadRequestInvalidRange (string): error name for invalid range
const errorBadRequestInvalidRange = "InvalidRange"

// errorInvalidAuthentication (string): error name for invalid AuthN
const errorInvalidAuthentication = "InvalidAuthentication"

// errorPermissionDenied (string): error for permission denied (invalid AuthZ)
const errorPermissionDenied = "PermissionDenied"

// errorNotFound (string): error name for not found
const errorNotFound = "NotFound"

// errorInternalServerError (string): error name for unspecified server errors
const errorInternalServerError = "InternalServerError"

/* Default Messages: default error message by error name */

// dfltMsgBadRequestUnsupportedFormat (string): default unsupported format message
const dfltMsgBadRequestUnsupportedFormat = "The requested file format is not supported by the server"

// dfltMsgBadRequestInvalidInput (string): default invalid input message
const dfltMsgBadRequestInvalidInput = "The request parameters do not adhere to the specification"

// dfltMsgBadRequestInvalidRange (string): default invalid range message
const dfltMsgBadRequestInvalidRange = "The requested range cannot be satisfied"

// dfltMsgInvalidAuthentication (string): default AuthN error message
const dfltMsgInvalidAuthentication = "Authorization provided is invalid"

// dfltMsgPermissionDenied (string): default AuthZ error message
const dfltMsgPermissionDenied = "Authorization is required to access the resource"

// dfltMsgNotFound (string): default not found error message
const dfltMsgNotFound = "The resource requested was not found"

// dfltMsgInternalServerError (string): default message for unspecified errors
const dfltMsgInternalServerError = "Internal server error"

// errorInfoMap (map[string]map[string]string) maps error name to status code
// and default message
var errorInfoMap = map[string]map[string]string{
	errorBadRequestUnsupportedFormat: {
		"code":    strconv.Itoa(codeBadRequest),
		"dfltMsg": dfltMsgBadRequestUnsupportedFormat,
	},
	errorBadRequestInvalidInput: {
		"code":    strconv.Itoa(codeBadRequest),
		"dfltMsg": dfltMsgBadRequestInvalidInput,
	},
	errorBadRequestInvalidRange: {
		"code":    strconv.Itoa(codeBadRequest),
		"dfltMsg": dfltMsgBadRequestInvalidRange,
	},
	errorInvalidAuthentication: {
		"code":    strconv.Itoa(codeInvalidAuthentication),
		"dfltMsg": dfltMsgInvalidAuthentication,
	},
	errorPermissionDenied: {
		"code":    strconv.Itoa(codePermissionDenied),
		"dfltMsg": dfltMsgPermissionDenied,
	},
	errorNotFound: {
		"code":    strconv.Itoa(codeNotFound),
		"dfltMsg": dfltMsgNotFound,
	},
	errorInternalServerError: {
		"code":    strconv.Itoa(codeInternalServerError),
		"dfltMsg": dfltMsgInternalServerError,
	},
}
