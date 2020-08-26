// Package main contains the main method/entrypoint
//
// Module main.go contains the main method/entrypoint
package main

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsserver"
)

// main program entrypoint
func main() {

	// load configuration object
	htsconfig.GetConfig()
	configLoadError := htsconfig.GetConfigLoadError()
	if configLoadError != nil {
		panic(configLoadError.Error())
	}

	// load server routes
	router, err := htsserver.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}

	// start server
	port := htsconfig.GetPort()
	fmt.Printf("Server started on port %s!\n", port)
	http.ListenAndServe(":"+port, router)
}
