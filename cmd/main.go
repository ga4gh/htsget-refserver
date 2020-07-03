package main

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/config"
	"github.com/ga4gh/htsget-refserver/internal/htsgetserver"
)

func main() {
	router, err := htsgetserver.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}
	port := config.GetConfigProp("port")
	fmt.Printf("Server started on port %s!\n", port)
	http.ListenAndServe(":"+port, router)
}
