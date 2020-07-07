package main

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsgetserver"
)

func main() {

	// load configuration object
	htsgetconfig.LoadAndValidateConfig()
	configLoadError := htsgetconfig.GetConfigLoadError()
	if configLoadError != nil {
		panic(configLoadError.Error())
	}
	// load server routes
	router, err := htsgetserver.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}
	// start server
	port := htsgetconfig.GetPort()
	fmt.Printf("Server started on port %s!\n", port)
	http.ListenAndServe(":"+port, router)
}
