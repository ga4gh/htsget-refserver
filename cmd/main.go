// Package main contains the main method/entrypoint
//
// Module main.go contains the main method/entrypoint
package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsserver"
	log "github.com/sirupsen/logrus"
)

// main program entrypoint
func main() {

	TLSCert := flag.String("cert", "/certs/server.crt", "TLS Server certificate")
	TLSKey := flag.String("key", "/certs/server.key", "TLS encryption key")
	NoTLS := flag.Bool("notls", false, "Run server without TLS")

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
	server := &http.Server{
		Addr:              ":" + htsconfig.GetPort(),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           logRequest(http.DefaultServeMux),
	}

	if *NoTLS {
		log.Infof("Insecure HTTP Server started at %s", server.Addr)
		log.Fatal(server.ListenAndServe())
	} else {
		log.Infof("HTTPS Server started at %s", server.Addr)
		log.Fatal(server.ListenAndServeTLS(*TLSCert, *TLSKey))
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
