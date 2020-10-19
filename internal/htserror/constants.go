// Package htserror write HTTP errors with various codes, response bodies
//
// Module constants contains HTTP Error constants (status codes, error names,
// default messages)
package htserror

import (
	"net/http"
	"strconv"
)

/* Error Status Codes: HTTP status code constants by error */

// codeBadRequest status code for Bad Request errors
const codeBadRequest = http.StatusBadRequest

// codeInvalidAuthentication status code for AuthN errors
const codeInvalidAuthentication = http.StatusUnauthorized

// codePermissionDenied status code for AuthZ errors
const codePermissionDenied = http.StatusForbidden

// codeNotFound status code for Not Found errors
const codeNotFound = http.StatusNotFound

// codeInternalServerError status code for unspecified server-side error
const codeInternalServerError = http.StatusInternalServerError

/* Error Names: htsget canonical error names */

// errorBadRequestUnsupportedFormat error name for unsupported format
const errorBadRequestUnsupportedFormat = "UnsupportedFormat"

// errorBadRequestInvalidInput error name for invalid input
const errorBadRequestInvalidInput = "InvalidInput"

// errorBadRequestInvalidRange error name for invalid range
const errorBadRequestInvalidRange = "InvalidRange"

// errorInvalidAuthentication error name for invalid AuthN
const errorInvalidAuthentication = "InvalidAuthentication"

// errorPermissionDenied error for permission denied (invalid AuthZ)
const errorPermissionDenied = "PermissionDenied"

// errorNotFound error name for not found
const errorNotFound = "NotFound"

// errorInternalServerError error name for unspecified server errors
const errorInternalServerError = "InternalServerError"

/* Default Messages: default error message by error name */

// dfltMsgBadRequestUnsupportedFormat default unsupported format message
const dfltMsgBadRequestUnsupportedFormat = "The requested file format is not supported by the server"

// dfltMsgBadRequestInvalidInput default invalid input message
const dfltMsgBadRequestInvalidInput = "The request parameters do not adhere to the specification"

// dfltMsgBadRequestInvalidRange default invalid range message
const dfltMsgBadRequestInvalidRange = "The requested range cannot be satisfied"

// dfltMsgInvalidAuthentication default AuthN error message
const dfltMsgInvalidAuthentication = "Authorization provided is invalid"

// dfltMsgPermissionDenied default AuthZ error message
const dfltMsgPermissionDenied = "Authorization is required to access the resource"

// dfltMsgNotFound default not found error message
const dfltMsgNotFound = "The resource requested was not found"

// dfltMsgInternalServerError default message for unspecified errors
const dfltMsgInternalServerError = "Internal server error"

// errorInfoMap maps error name to status code and default message
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
