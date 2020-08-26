package htsserviceinfo

import (
	"sync"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

type ServiceType struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}

type Organization struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ServiceInfo struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Type             *ServiceType  `json:"type"`
	Description      string        `json:"description"`
	Organization     *Organization `json:"organization"`
	ContactURL       string        `json:"contactUrl"`
	documentationURL string        `json:"documentationUrl"`
	CreatedAt        string        `json:"createdAt"`
	UpdatedAt        string        `json:"updatedAt"`
	Environment      string        `json:"environment"`
	Version          string        `json:"version"`
}

var readsServiceInfoSingleton *ServiceInfo
var variantsServiceInfoSingleton *ServiceInfo
var serviceInfoSingletonMap = map[htsconstants.APIEndpoint]*ServiceInfo{
	htsconstants.APIEndpointReadsServiceInfo:    readsServiceInfoSingleton,
	htsconstants.APIEndpointVariantsServiceInfo: variantsServiceInfoSingleton,
}

var readsServiceInfoLoaded sync.Once
var variantsServiceInfoLoaded sync.Once
var serviceInfoLoadedMap = map[htsconstants.APIEndpoint]*sync.Once{
	htsconstants.APIEndpointReadsServiceInfo:    &readsServiceInfoLoaded,
	htsconstants.APIEndpointVariantsServiceInfo: &variantsServiceInfoLoaded,
}

func loadServiceInfo(endpoint htsconstants.APIEndpoint) {
	newServiceInfo := new(ServiceInfo)
	newServiceType := new(ServiceType)
	newServiceType.Group = "org.ga4gh"
	newServiceType.Artifact = "htsget"
	newServiceType.Version = "1.2.0"
	newServiceInfo.Type = newServiceType
	serviceInfoSingletonMap[endpoint] = newServiceInfo
}

func GetServiceInfo(endpoint htsconstants.APIEndpoint) *ServiceInfo {
	loadedCheck := serviceInfoLoadedMap[endpoint]
	loadedCheck.Do(func() {
		loadServiceInfo(endpoint)
	})
	return serviceInfoSingletonMap[endpoint]
}
