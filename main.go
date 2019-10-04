package main

import (
	"net/http"

	"github.com/david-xliu/htsget-refserver/handler"
)

func main() {
	router := handler.SetRouter()
	http.ListenAndServe(":3000", router)
}
