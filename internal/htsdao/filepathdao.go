package htsdao

import (
	"math"
	"net/http"
	"os"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsticket"
)

type FilePathDao struct {
	id       string
	filePath string
}

func NewFilePathDao(id string, filePath string) *FilePathDao {
	dao := new(FilePathDao)
	dao.id = id
	dao.filePath = filePath
	return dao
}

func (dao *FilePathDao) GetContentLength(request *http.Request) int64 {
	fileInfo, _ := os.Stat(dao.filePath)
	return fileInfo.Size()
}

func (dao *FilePathDao) constructByteRangeURL(start int64, end int64, request *http.Request) *htsticket.URL {
	host := htsconfig.GetHost()
	path := host + htsconstants.FileByteRangeURLPath
	headers := htsticket.NewHeaders()
	headers.SetRangeHeader(start, end)
	headers.SetFilePathHeader(dao.filePath)
	url := htsticket.NewURL()
	url.SetURL(path)
	url.SetHeaders(headers)
	url.Headers.Authorization = request.Header.Get("Authorization")

	return url
}

func (dao *FilePathDao) GetByteRangeUrls(request *http.Request) []*htsticket.URL {
	numBytes := dao.GetContentLength(request)
	blockSize := htsconstants.SingleBlockByteSize
	var start, end int64 = 0, 0
	numBlocks := int(math.Ceil(float64(numBytes) / float64(blockSize)))
	urls := []*htsticket.URL{}
	for i := 1; i <= numBlocks; i++ {
		end = start + blockSize - 1
		if end >= numBytes {
			end = numBytes - 1
		}
		url := dao.constructByteRangeURL(start, end, request)
		start = end + 1
		urls = append(urls, url)
	}
	return urls
}

func (dao *FilePathDao) String() string {
	return "FilePathDao id=" + dao.id + ", filePath=" + dao.filePath
}
