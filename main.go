package main

import (
	"fmt"
	"net/http"

	"github.com/david-xliu/htsget-refserver/handler"
)

func main() {
	router := handler.SetRouter()
	fmt.Println("Server started on port 3000!")
	http.ListenAndServe(":3000", router)
}
