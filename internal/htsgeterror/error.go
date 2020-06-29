// Package htsgeterror provides operations for writing various HTTP errors,
// including different status codes and error response bodies
package htsgeterror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// htsgetError contains attributes to write an error as an HTTP response in Go.
// contains the code, and a container for the response body
type htsgetError struct {
	Code   int
	Htsget errorContainer `json:"htsget"`
}

// errorContainer contains attributes for an htsget error response body,
// including the error name and a helper message
type errorContainer struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Error displays htsgetError attributes as a string
func (err *htsgetError) Error() string {
	return fmt.Sprint(err.Htsget.Error + ": " + err.Htsget.Message)
}

// instantiates an htsgetError object with a specific code, name, and message
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

// writes a passed htsget error as json to the HTTP ResponseWriter
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

// template function for more specific htsget error types. writes an htsget
// error to the ResponseWriter. uses the default message in the response body
// if no message is passed
func htsgetErrorTemplate(writer http.ResponseWriter, errorName string, msgPtr *string) {
	code, codeParseErr := strconv.Atoi(errorInfoMap[errorName]["code"])
	if codeParseErr != nil {

	}
	msg := errorInfoMap[errorName]["dfltMsg"]
	if msgPtr != nil {
		msg = *msgPtr
	}
	writeHtsgetErrorToHTTPError(writer, code, errorName, msg)
}

// UnsupportedFormat writes an UnsupportedFormat error to the HTTP ResponseWriter
func UnsupportedFormat(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestUnsupportedFormat, msgPtr)
}

// InvalidInput writes an InvalidInput error to the HTTP ResponseWriter
func InvalidInput(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidInput, msgPtr)
}

// InvalidRange writes an InvalidRange error to the HTTP ResponseWriter
func InvalidRange(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidRange, msgPtr)
}

// InvalidAuthentication writes an InvalidAuthentication error to the HTTP ResponseWriter
func InvalidAuthentication(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInvalidAuthentication, msgPtr)
}

// PermissionDenied writes a PermissionDenied error to the HTTP ResponseWriter
func PermissionDenied(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorPermissionDenied, msgPtr)
}

// NotFound writes a NotFound error to the HTTP ResponseWriter
func NotFound(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorNotFound, msgPtr)
}

// InternalServerError writes an InternalServerError error to the HTTP ResponseWriter
func InternalServerError(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInternalServerError, msgPtr)
}
