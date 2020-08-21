package htsdao

import "github.com/ga4gh/htsget-refserver/internal/htsticket"

type DataAccessObject interface {
	GetContentLength() int64
	GetByteRangeUrls() []*htsticket.URL
	String() string
}
