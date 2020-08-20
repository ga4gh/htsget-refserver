package htsgetdao

import (
	"math"
	"os"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig/htsgetconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetticket"
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

func (dao *FilePathDao) GetContentLength() int64 {
	fileInfo, _ := os.Stat(dao.filePath)
	return fileInfo.Size()
}

func (dao *FilePathDao) constructByteRangeURL(start int64, end int64) *htsgetticket.URL {
	host := htsgetconfig.GetHost()
	path := host + htsgetconstants.FileByteRangeURLPath
	headers := htsgetticket.NewHeaders()
	headers.SetRangeHeader(start, end)
	headers.SetFilePathHeader(dao.filePath)
	url := htsgetticket.NewURL()
	url.SetURL(path)
	url.SetHeaders(headers)
	return url
}

func (dao *FilePathDao) GetByteRangeUrls() []*htsgetticket.URL {
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
		url := dao.constructByteRangeURL(start, end)
		start = end + 1
		urls = append(urls, url)
	}
	return urls
}

func (dao *FilePathDao) String() string {
	return "FilePathDao id=" + dao.id + ", filePath=" + dao.filePath
}
