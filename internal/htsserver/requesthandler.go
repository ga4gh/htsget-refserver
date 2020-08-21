package htsserver

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

type requestHandler struct {
	method          htsconstants.HttpMethod
	endpoint        htsconstants.ServerEndpoint
	handlerFunction func(handler *requestHandler)
	Writer          http.ResponseWriter
	Request         *http.Request
	HtsReq          *htsrequest.HtsgetRequest
}

func newRequestHandler(method htsconstants.HttpMethod, endpoint htsconstants.ServerEndpoint, handlerFunction func(handler *requestHandler)) *requestHandler {
	reqHandler := new(requestHandler)
	reqHandler.method = method
	reqHandler.endpoint = endpoint
	reqHandler.handlerFunction = handlerFunction
	return reqHandler
}

func (reqHandler *requestHandler) stage(writer http.ResponseWriter, request *http.Request) error {
	htsgetReq, err := htsrequest.SetAllParameters(reqHandler.method, reqHandler.endpoint, writer, request)

	if err != nil {
		return err
	}

	reqHandler.Writer = writer
	reqHandler.Request = request
	reqHandler.HtsReq = htsgetReq
	return nil
}

func (reqHandler *requestHandler) execute() {
	reqHandler.handlerFunction(reqHandler)
}

func (reqHandler *requestHandler) handleRequest(writer http.ResponseWriter, request *http.Request) error {
	stagingErr := reqHandler.stage(writer, request)
	if stagingErr != nil {
		return stagingErr
	}
	reqHandler.execute()
	return nil
}
