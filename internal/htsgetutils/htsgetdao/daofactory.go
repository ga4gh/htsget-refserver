package htsgetdao

import (
	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
)

func GetMatchingDao(id string, registry *htsgetconfig.DataSourceRegistry) (DataAccessObject, error) {
	path, err := registry.GetMatchingPath(id)
	if err != nil {
		return nil, err
	}
	if htsgetutils.IsValidUrl(path) {
		return NewURLDao(id, path), nil
	} else {
		return NewFilePathDao(id, path), nil
	}
}

func GetReadsDaoForID(id string) (DataAccessObject, error) {
	return GetMatchingDao(id, htsgetconfig.GetReadsDataSourceRegistry())
}
