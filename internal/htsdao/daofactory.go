package htsdao

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

func GetMatchingDao(id string, registry *htsconfig.DataSourceRegistry) (DataAccessObject, error) {
	path, err := registry.GetMatchingPath(id)
	if err != nil {
		return nil, err
	}
	if htsutils.IsValidUrl(path) {
		return NewURLDao(id, path), nil
	} else {
		return NewFilePathDao(id, path), nil
	}
}

func GetReadsDaoForID(id string) (DataAccessObject, error) {
	return GetMatchingDao(id, htsconfig.GetReadsDataSourceRegistry())
}

func GetVariantsDaoForID(id string) (DataAccessObject, error) {
	return GetMatchingDao(id, htsconfig.GetVariantsDataSourceRegistry())
}
