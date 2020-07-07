package htsgetdao

import "github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetticket"

type DataAccessObject interface {
	GetContentLength() int64
	GetByteRangeUrls() []*htsgetticket.URL
	String() string
}
