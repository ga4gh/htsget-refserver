package htsserver

import (
	"encoding/json"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

func serviceInfoRequestHandler(handler *requestHandler) {
	serviceInfo := handler.HtsReq.GetServiceInfo()
	writer := handler.Writer
	writer.Header().Set(htsconstants.ContentTypeHeader.String(), htsconstants.ContentTypeHeaderHtsgetJSON.String())
	json.NewEncoder(writer).Encode(serviceInfo)
}
