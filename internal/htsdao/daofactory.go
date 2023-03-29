package htsdao

import (
	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
	log "github.com/sirupsen/logrus"
)

func getMatchingDao(id string, registry *htsconfig.DataSourceRegistry) (DataAccessObject, error) {
	path, err := registry.GetMatchingPath(id)
	if err != nil {
		log.Debugf("error in GetMatchingPath: %v", err)
		return nil, err
	}
	if htsutils.IsValidURL(path) {
		return NewURLDao(id, path), nil
	}
	return NewFilePathDao(id, path), nil
}

func GetDao(req *htsrequest.HtsgetRequest) (DataAccessObject, error) {
	registry := req.GetDataSourceRegistry()
	return getMatchingDao(req.GetID(), registry)
}
