package htsgetdao

import (
	"math"
	"net/http"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig/htsgetconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetticket"
)

type URLDao struct {
	id  string
	url string
}

func NewURLDao(id string, url string) *URLDao {
	dao := new(URLDao)
	dao.id = id
	dao.url = url
	return dao
}

func (dao *URLDao) GetContentLength() int64 {
	res, _ := http.Head(dao.url)
	return res.ContentLength
}

func (dao *URLDao) GetByteRangeUrls() []*htsgetticket.URL {

	numBytes := dao.GetContentLength()
	blockSize := htsgetconstants.SingleBlockByteSize
	var start, end int64 = 0, 0
	numBlocks := int(math.Ceil(float64(numBytes) / float64(blockSize)))
	urls := []*htsgetticket.URL{}
	for i := 1; i <= numBlocks; i++ {
		end = start + blockSize - 1
		if end >= numBytes {
			end = numBytes - 1
		}
		headers := htsgetticket.NewHeaders()
		headers.SetRangeHeader(start, end)
		url := htsgetticket.NewURL()
		url.SetURL(dao.url)
		url.SetHeaders(headers)
		start = end + 1
		urls = append(urls, url)
	}
	return urls
}

func (dao *URLDao) String() string {
	return "URLDao id=" + dao.id + ", url=" + dao.url
}
