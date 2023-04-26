// Package main contains the main method/entrypoint
//
// Module main.go contains the main method/entrypoint
package main

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsserver"
	log "github.com/sirupsen/logrus"
)

// main program entrypoint
func main() {

	log.SetLevel(log.InfoLevel)
	// load configuration object
	htsconfig.GetConfig()
	configLoadError := htsconfig.GetConfigLoadError()
	if configLoadError != nil {
		log.Errorf("error from getConfigLoadError: %v", configLoadError)
		panic(configLoadError.Error())
	}

	// load server routes
	router, err := htsserver.SetRouter()
	if err != nil {
		log.Errorf("error setting up router: %v", err)
		panic("Problem setting up server.")
	}
	http.Handle("/", router)

	// start server
	port := htsconfig.GetPort()
	fmt.Printf("Server started on port %s!\n", port)
	http.ListenAndServe(":"+port, logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
