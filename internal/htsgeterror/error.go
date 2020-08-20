// Package htsgeterror provides operations for writing various HTTP errors,
// including different status codes and error response bodies
//
// Module error.go contains operations for constructing htsget errors, and
// writing them to the HTTP response
package htsgeterror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// htsgetError contains attributes to write an error as an HTTP response
//
// Attributes
//	Code (int): HTTP status code
//	Htsget (errorContainer): contains error name and message
type htsgetError struct {
	Code   int
	Htsget errorContainer `json:"htsget"`
}

// errorContainer contains attributes for an htsget error response body
//
// Attributes
//	Error (string): htsget error name (according to htsget specification)
//	Message (string): client help message, explaining why error was encountered
type errorContainer struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Error displays htsgetError attributes as a string
//
// Type: htsgetError
// Returns
//	(string): string representation of the error (name and message)
func (err *htsgetError) Error() string {
	return fmt.Sprint(err.Htsget.Error + ": " + err.Htsget.Message)
}

// newHtsgetError instantiates a new htsgetError instance
//
// Arguments
//	code (int): HTTP status code
//	err (string): error name
//	message (string): error message
// Returns
//	(htsgetError): error object with attributes set based on passed arguments
func newHtsgetError(code int, err string, message string) *htsgetError {
	htsgetError := &htsgetError{
		Code: code,
		Htsget: errorContainer{
			err,
			message,
		},
	}
	return htsgetError
}

// writeHTTPError writes an htsgetError code and body to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	err (error): any error object, including an htsgetError
func writeHTTPError(writer http.ResponseWriter, err error) {
	if err, ok := err.(*htsgetError); ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(err.Code)
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"htsget": err.Htsget,
		})
		return
	}
	http.Error(writer, err.Error(), 500)
}

// instantiates an htsgetError based on passed values, and writes it to the
// HTTP ResponseWriter
func writeHtsgetErrorToHTTPError(writer http.ResponseWriter, code int, err string, message string) {
	htsgetError := newHtsgetError(code, err, message)
	writeHTTPError(writer, htsgetError)
}

// htsgetErrorTemplate is a template function for more specific htsget error
// types. writes an htsget error to the ResponseWriter. uses the default message
// in the response body if no message is passed
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	errorName (string): htsget-specific error name
//	msgPtr (*string): pointer to a message string, if nil, use the default message
func htsgetErrorTemplate(writer http.ResponseWriter, errorName string, msgPtr *string) {
	code, codeParseErr := strconv.Atoi(errorInfoMap[errorName]["code"])
	if codeParseErr != nil {
	}

	// use the default message if the message pointer is nil
	msg := errorInfoMap[errorName]["dfltMsg"]
	if msgPtr != nil {
		msg = *msgPtr
	}
	writeHtsgetErrorToHTTPError(writer, code, errorName, msg)
}

// UnsupportedFormat writes an UnsupportedFormat error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
// msgPtr (*string): pointer to message string, if nil, use the default message
func UnsupportedFormat(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestUnsupportedFormat, msgPtr)
}

// InvalidInput writes an InvalidInput error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	msgPtr (*string): pointer to message string, if nil, use the default message
func InvalidInput(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidInput, msgPtr)
}

// InvalidRange writes an InvalidRange error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	msgPtr (*string): pointer to message string, if nil, use the default message
func InvalidRange(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidRange, msgPtr)
}

// InvalidAuthentication writes an InvalidAuthentication error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	msgPtr (*string): pointer to message string, if nil, use the default message
func InvalidAuthentication(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInvalidAuthentication, msgPtr)
}

// PermissionDenied writes a PermissionDenied error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	msgPtr (*string): pointer to message string, if nil, use the default message
func PermissionDenied(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorPermissionDenied, msgPtr)
}

// NotFound writes a NotFound error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
// 	msgPtr (*string): pointer to message string, if nil, use the default message
func NotFound(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorNotFound, msgPtr)
}

// InternalServerError writes an InternalServerError error to the HTTP ResponseWriter
//
// Arguments
//	writer (http.ResponseWriter): HTTP response writer, sets response code and body
//	msgPtr (*string): pointer to message string, if nil, use the default message
func InternalServerError(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInternalServerError, msgPtr)
}
