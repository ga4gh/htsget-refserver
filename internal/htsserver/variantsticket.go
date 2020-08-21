package htsserver

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htserror"

	"github.com/ga4gh/htsget-refserver/internal/htsdao"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

func getVariantsTicket(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("A")

	params := request.URL.Query()
	fmt.Println("B")
	// host := htsgetconfig.GetHost()
	htsgetReq, err := htsrequest.VariantsTicketEndpointSetAllParameters(request, writer, params)
	fmt.Println("C")

	if err != nil {
		return
	}

	fmt.Println("D")
	//dao, err := htsdao.GetVariantsDaoForID(htsgetReq.ID())
	_, err = htsdao.GetVariantsDaoForID(htsgetReq.ID())
	fmt.Println("E")

	if err != nil {
		msg := "Could not determine data source path/url from request id"
		htserror.InternalServerError(writer, &msg)
	}
	fmt.Println("F")
}
