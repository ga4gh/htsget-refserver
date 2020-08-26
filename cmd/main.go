// Package main contains the main method/entrypoint
//
// Module main.go contains the main method/entrypoint
package main

import (
	"fmt"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
)

// main program entrypoint
func main() {
	// load configuration object
	fmt.Println("START")
	fmt.Println(htsconfig.Port())
	fmt.Println(htsconfig.Host())
	fmt.Println("Is my reads endpoint loaded?")
	fmt.Println(htsconfig.IsEndpointEnabled(htsconstants.APIEndpointReadsTicket))

	/*
		htsconfig.LoadAndValidateConfig()
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
	*/
}
