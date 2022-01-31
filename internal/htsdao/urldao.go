package htsdao

import (
	"math"
	"net/http"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
	"github.com/ga4gh/htsget-refserver/internal/awsutils"
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
	if strings.HasPrefix(dao.url, awsutils.S3Proto) {
		contentLength, _ := awsutils.HeadS3Object(awsutils.S3Dto{
			ObjPath: dao.url,
		})
		return contentLength
	}
	res, _ := http.Head(dao.url)
	return res.ContentLength
}

func (dao *URLDao) GetByteRangeUrls() []*htsticket.URL {

	numBytes := dao.GetContentLength()
	blockSize := htsconstants.SingleBlockByteSize
	var start, end int64 = 0, 0
	numBlocks := int(math.Ceil(float64(numBytes) / float64(blockSize)))
	urls := []*htsticket.URL{}
	for i := 1; i <= numBlocks; i++ {
		end = start + blockSize - 1
		if end >= numBytes {
			end = numBytes - 1
		}
		headers := htsticket.NewHeaders()
		headers.SetRangeHeader(start, end)
		url := htsticket.NewURL()
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
