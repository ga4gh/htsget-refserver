package main

import (
	"fmt"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsgetserver"
)

func main() {
	router, err := htsgetserver.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}
	port := htsgetconfig.GetPort()
	fmt.Printf("Server started on port %s!\n", port)
	http.ListenAndServe(":"+port, router)
}
