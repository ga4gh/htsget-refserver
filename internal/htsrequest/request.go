// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module request.go defines structs and operations for a mature htsget
// request, which holds all htsget-related parameters and has insight into
// what information was requested by the client
package htsrequest

import (
	"net/url"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsserviceinfo"

	"github.com/ga4gh/htsget-refserver/internal/htsutils"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
)

// HtsgetRequest contains htsget-related parameters
//
// Attributes
//	ScalarParams (map[string]string): map holding scalar parameter values
//	ListParams (map[string][]string): map holding list parameter values
type HtsgetRequest struct {
	endpoint     htsconstants.APIEndpoint
	ScalarParams map[string]string
	ListParams   map[string][]string
}

// NewHtsgetRequest instantiates a new HtsgetRequest struct instance
//
// Returns
//	(*HtsgetRequest): new HtsgetRequest instance
func NewHtsgetRequest() *HtsgetRequest {
	htsgetReq := new(HtsgetRequest)
	htsgetReq.ScalarParams = make(map[string]string)
	htsgetReq.ListParams = make(map[string][]string)
	return htsgetReq
}

func (htsgetReq *HtsgetRequest) SetEndpoint(endpoint htsconstants.APIEndpoint) {
	htsgetReq.endpoint = endpoint
}

func (htsgetReq *HtsgetRequest) GetEndpoint() htsconstants.APIEndpoint {
	return htsgetReq.endpoint
}

// AddScalarParam adds a key-value pair to HtsgetRequest scalar parameter map
//
// Type: HtsgetRequest
// Arguments
//	key (string): parameter name
//	value (string): parameter value
func (htsgetReq *HtsgetRequest) AddScalarParam(key string, value string) {
	htsgetReq.ScalarParams[key] = value
}

// AddListParam adds a key-value pair to HtsgetRequest list parameter map
//
// Type: HtsgetRequest
// Arguments
//	key (string): parameter name
//	value ([]string): parameter value
func (htsgetReq *HtsgetRequest) AddListParam(key string, value []string) {
	htsgetReq.ListParams[key] = value
}

// get retrieves a value from the scalar parameter map by its key
//
// Type: HtsgetRequest
// Arguments
//	key (string): parameter name
// Returns
//	(string): parameter value under the given key/name
func (htsgetReq *HtsgetRequest) get(key string) string {
	return htsgetReq.ScalarParams[key]
}

// getList retrieves a value from the list parameter map by its key
//
// Type: HtsgetRequest
// Arguments
//	key (string): parameter name
// Returns
//	(string): parameter value under the given key/name
func (htsgetReq *HtsgetRequest) getList(key string) []string {
	return htsgetReq.ListParams[key]
}

// ID gets value of 'id' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'id'
func (htsgetReq *HtsgetRequest) ID() string {
	return htsgetReq.get("id")
}

// Format gets value of 'format' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'format'
func (htsgetReq *HtsgetRequest) Format() string {
	return htsgetReq.get("format")
}

// Class gets value of 'class' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'class'
func (htsgetReq *HtsgetRequest) Class() string {
	return htsgetReq.get("class")
}

// ReferenceName gets value of 'referenceName' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'referenceName'
func (htsgetReq *HtsgetRequest) ReferenceName() string {
	return htsgetReq.get("referenceName")
}

// Start gets value of 'start' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'start'
func (htsgetReq *HtsgetRequest) Start() string {
	return htsgetReq.get("start")
}

// End gets value of 'end' param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'end'
func (htsgetReq *HtsgetRequest) End() string {
	return htsgetReq.get("end")
}

// HtsgetBlockClass gets value of 'HtsgetBlockClass' header param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'HtsgetBlockClass'
func (htsgetReq *HtsgetRequest) HtsgetBlockClass() string {
	return htsgetReq.get("HtsgetBlockClass")
}

// HtsgetBlockID gets value of 'HtsgetBlockId' header param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'HtsgetBlockId'
func (htsgetReq *HtsgetRequest) HtsgetBlockID() string {
	return htsgetReq.get("HtsgetBlockId")
}

// HtsgetNumBlocks gets value of 'HtsgetNumBlocks' header param
//
// Type: HtsgetRequest
// Returns
//	(string): value of 'HtsgetNumBlocks'
func (htsgetReq *HtsgetRequest) HtsgetNumBlocks() string {
	return htsgetReq.get("HtsgetNumBlocks")
}

func (htsgetReq *HtsgetRequest) HtsgetFilePath() string {
	return htsgetReq.get("HtsgetFilePath")
}

func (htsgetReq *HtsgetRequest) Range() string {
	return htsgetReq.get("Range")
}

// Fields gets value of 'fields' param
//
// Type: HtsgetRequest
// Returns
//	([]string): value of 'fields'
func (htsgetReq *HtsgetRequest) Fields() []string {
	return htsgetReq.getList("fields")
}

// Tags gets value of 'tags' param
//
// Type: HtsgetRequest
// Returns
//	([]string): value of 'tags'
func (htsgetReq *HtsgetRequest) Tags() []string {
	return htsgetReq.getList("tags")
}

// NoTags gets value of 'notags' param
//
// Type: HtsgetRequest
// Returns
//	([]string): value of 'notags'
func (htsgetReq *HtsgetRequest) NoTags() []string {
	return htsgetReq.getList("notags")
}

// isDefaultScalar checks if a scalar parameter value matches the default,
// unspecfied value, thereby indicating that the parameter was not specified
// in the HTTP request
//
// Type: HtsgetRequest
// Arguments
//	key (string): the parameter name to check
// Returns
//	(bool): true if the parameter value matches the default, false if not
func (htsgetReq *HtsgetRequest) isDefaultScalar(key string) bool {
	return htsgetReq.get(key) == defaultScalarParameterValues[key]
}

// isDefaultScalar checks if a list parameter value matches the default,
// unspecfied value, thereby indicating that the parameter was not specified
// in the HTTP request
//
// Type: HtsgetRequest
// Arguments
//	key (string): the parameter name to check
// Returns
//	(bool): true if the parameter value matches the default, false if not
func (htsgetReq *HtsgetRequest) isDefaultList(key string) bool {
	list := htsgetReq.getList(key)
	defaultList := defaultListParameterValues[key]
	if len(list) != len(defaultList) {
		return false
	}
	for i := 0; i < len(list); i++ {
		if list[i] != defaultList[i] {
			return false
		}
	}
	return true
}

// HeaderOnlyRequested checks if the client request is only for the header
//
// Type: HtsgetRequest
// Returns
//	(bool): true if only the header was requested
func (htsgetReq *HtsgetRequest) HeaderOnlyRequested() bool {
	return htsgetReq.Class() == htsconstants.ClassHeader
}

// UnplacedUnmappedReadsRequested checks if the client request is for unplaced,
// unmapped reads
//
// Type: HtsgetRequest
// Returns
//	(bool): true if unplaced, unmapped reads were requested
func (htsgetReq *HtsgetRequest) UnplacedUnmappedReadsRequested() bool {
	return htsgetReq.ReferenceName() == "*"
}

func (htsgetReq *HtsgetRequest) ReferenceNameRequested() bool {
	return !htsgetReq.isDefaultScalar("referenceName")
}

func (htsgetReq *HtsgetRequest) StartRequested() bool {
	return !htsgetReq.isDefaultScalar("start")
}

func (htsgetReq *HtsgetRequest) EndRequested() bool {
	return !htsgetReq.isDefaultScalar("end")
}

// AllRegionsRequested checks if the client request is for all chromosomal
// regions in the file (ie. referenceName not specified)
//
// Type: HtsgetRequest
// Returns
//	(bool): true if all chromosomal regions requested
func (htsgetReq *HtsgetRequest) AllRegionsRequested() bool {
	return htsgetReq.isDefaultScalar("referenceName")
}

// AllFieldsRequested checks if all fields were requested by the client. all
// fields are requested if the client does not specify the 'fields' parameter
//
// Type: HtsgetRequest
// Returns
//	(bool): true if all fields were requested by client, false if not
func (htsgetReq *HtsgetRequest) AllFieldsRequested() bool {
	return htsgetReq.isDefaultList("fields")
}

// TagsNotSpecified checks if the tags parameter was not provided in the HTTP
// request
//
// Type: HtsgetRequest
// Returns
// (bool): true if the tags parameter was not specified in request, false if not
func (htsgetReq *HtsgetRequest) TagsNotSpecified() bool {
	return htsgetReq.isDefaultList("tags")
}

// NoTagsNotSpecified checks if the notags parameter was not provided in the
// HTTP request
//
// Type: HtsgetRequest
// Returns
//	(bool): true if the notags parameter was not specified in request, false if not
func (htsgetReq *HtsgetRequest) NoTagsNotSpecified() bool {
	return htsgetReq.isDefaultList("notags")
}

// AllTagsRequested checks if all tags were requested by the client. all tags
// are requested if the request specifies neither the 'tags' or 'notags' parameters
//
// Type: HtsgetRequest
// Returns
//	(bool): true if neither tags nor notags was specified, thereby requesting all tags
func (htsgetReq *HtsgetRequest) AllTagsRequested() bool {
	return htsgetReq.TagsNotSpecified() && htsgetReq.NoTagsNotSpecified()
}

// ConstructDataEndpointURL for a given htsget request object, return the url
// that will redirect the client to the correct data download endpoint with
// all necessary parameters and headers provided
func (htsgetReq *HtsgetRequest) ConstructDataEndpointURL() (*url.URL, error) {
	host := htsconfig.GetHost()
	dataEndpointPath := htsgetReq.GetEndpoint().DataEndpointPath()
	dataEndpoint, err := url.Parse(htsutils.RemoveTrailingSlash(host) + dataEndpointPath + htsgetReq.ID())
	if err != nil {
		return nil, err
	}

	// add query params
	query := dataEndpoint.Query()
	if htsgetReq.HeaderOnlyRequested() {
		query.Set("class", htsgetReq.Class())
	}
	if htsgetReq.ReferenceNameRequested() {
		query.Set("referenceName", htsgetReq.ReferenceName())
	}
	if htsgetReq.StartRequested() {
		query.Set("start", htsgetReq.Start())
	}
	if htsgetReq.EndRequested() {
		query.Set("end", htsgetReq.End())
	}
	if !htsgetReq.AllFieldsRequested() {
		f := strings.Join(htsgetReq.Fields(), ",")
		query.Set("fields", f)
	}
	if !htsgetReq.TagsNotSpecified() {
		t := strings.Join(htsgetReq.Tags(), ",")
		query.Set("tags", t)
	}
	if !htsgetReq.NoTagsNotSpecified() {
		nt := strings.Join(htsgetReq.NoTags(), ",")
		query.Set("notags", nt)
	}
	dataEndpoint.RawQuery = query.Encode()
	return dataEndpoint, nil
}

func (htsgetReq *HtsgetRequest) GetCorrespondingDataSourceRegistry() *htsconfig.DataSourceRegistry {
	reads := htsconfig.GetReadsDataSourceRegistry()
	variants := htsconfig.GetVariantsDataSourceRegistry()
	registries := map[htsconstants.APIEndpoint]*htsconfig.DataSourceRegistry{
		htsconstants.APIEndpointReadsTicket:         reads,
		htsconstants.APIEndpointReadsData:           reads,
		htsconstants.APIEndpointReadsServiceInfo:    nil,
		htsconstants.APIEndpointVariantsTicket:      variants,
		htsconstants.APIEndpointVariantsData:        variants,
		htsconstants.APIEndpointVariantsServiceInfo: nil,
		htsconstants.APIEndpointFileBytes:           nil,
	}
	return registries[htsgetReq.GetEndpoint()]
}

func (htsgetReq *HtsgetRequest) GetServiceInfo() *htsserviceinfo.ServiceInfo {
	return htsserviceinfo.GetServiceInfo(htsgetReq.GetEndpoint())
}
