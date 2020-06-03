package htsgetutils

import (
	"strings"
)

func AddTrailingSlash(url string) string {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}
