package main

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/david-xliu/htsget-refserver/handler"
)

func main() {
	router := handler.SetRouter()
	http.ListenAndServe(":3000", router)
}
