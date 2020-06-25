package htsgeterror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type htsgetError struct {
	Code   int
	Htsget errorContainer `json:"htsget"`
}

type errorContainer struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (err *htsgetError) Error() string {
	return fmt.Sprint(err.Htsget.Error + ": " + err.Htsget.Message)
}

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

func writeHtsgetErrorToHTTPError(writer http.ResponseWriter, code int, err string, message string) {
	htsgetError := newHtsgetError(code, err, message)
	writeHTTPError(writer, htsgetError)
}

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

func UnsupportedFormat(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestUnsupportedFormat, msgPtr)
}

func InvalidInput(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidInput, msgPtr)
}

func InvalidRange(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorBadRequestInvalidRange, msgPtr)
}

func InvalidAuthentication(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInvalidAuthentication, msgPtr)
}

func PermissionDenied(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorPermissionDenied, msgPtr)
}

func NotFound(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorNotFound, msgPtr)
}

func InternalServerError(writer http.ResponseWriter, msgPtr *string) {
	htsgetErrorTemplate(writer, errorInternalServerError, msgPtr)
}
