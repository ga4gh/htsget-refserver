// Package main contains the main method/entrypoint
//
// Module main.go contains the main method/entrypoint
package main

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	log "github.com/ga4gh/htsget-refserver/internal/htslog"
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

	// set up our global logging instance
	log.Setup(htsconfig.GetLogFile(), htsconfig.GetLogLevel())

	// load server routes
	router, err := htsserver.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}
	http.Handle("/", router)

	// start server
	port := htsconfig.GetPort()

	log.Info("Server started on port %s!\n", port)

	http.ListenAndServe(":"+port, nil)
}
