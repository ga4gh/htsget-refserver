// Package htserror write HTTP errors with various codes, response bodies
//
// Module error contains operations for constructing htsget errors, and
// writing them to the HTTP response
package htserror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// htsgetError contains attributes to write an error as an HTTP response,
// including response code and JSON body container
type htsgetError struct {
	Code   int
	Htsget errorContainer `json:"htsget"`
}

// errorContainer contains attributes for the main htsget error response body
type errorContainer struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Error displays htsgetError attributes as a string
func (err *htsgetError) Error() string {
	log.Debugf("some error: %v", err.Htsget.Error)
	return fmt.Sprint(err.Htsget.Error + ": " + err.Htsget.Message)
}

// newHtsgetError instantiates a new htsgetError instance
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
// in the response body if no message is passed. a message pointer can be passed
// to overall the default message
func htsgetErrorTemplate(writer http.ResponseWriter, errorName string, msgPtr *string) {
	code, _ := strconv.Atoi(errorInfoMap[errorName]["code"])

	// use the default message if the message pointer is nil
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
