package main

import (
	"fmt"
	"net/http"

	"github.com/david-xliu/htsget-refserver/internal/server"
)

func main() {
	router, err := server.SetRouter()
	if err != nil {
		panic("Problem setting up server.")
	}
	fmt.Println("Server started on port 3000!")
	http.ListenAndServe(":3000", router)
}
