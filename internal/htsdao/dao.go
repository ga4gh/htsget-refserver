package htsdao

import (
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

type DataAccessObject interface {
	GetContentLength(request *http.Request) int64
	GetByteRangeUrls(request *http.Request) []*htsticket.URL
	String() string
}
