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
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

// HtsgetRequest contains htsget-related parameters
//
// Attributes
//	ScalarParams (map[string]string): map holding scalar parameter values
//	ListParams (map[string][]string): map holding list parameter values
type HtsgetRequest struct {
	endpoint           htsconstants.APIEndpoint
	id                 string
	format             string
	class              string
	referenceName      string
	start              int
	end                int
	fields             []string
	tags               []string
	noTags             []string
	regions            []*Region
	htsgetBlockClass   string
	htsgetCurrentBlock string
	htsgetTotalBlocks  string
	htsgetFilePath     string
	htsgetRange        string
}

// NewHtsgetRequest instantiates a new HtsgetRequest struct instance
func NewHtsgetRequest() *HtsgetRequest {
	r := new(HtsgetRequest)
	r.SetRegions([]*Region{})
	return r
}

/* SETTERS AND GETTERS */

func (r *HtsgetRequest) SetEndpoint(endpoint htsconstants.APIEndpoint) {
	r.endpoint = endpoint
}

func (r *HtsgetRequest) GetEndpoint() htsconstants.APIEndpoint {
	return r.endpoint
}

func (r *HtsgetRequest) SetID(id string) {
	r.id = id
}

func (r *HtsgetRequest) GetID() string {
	return r.id
}

func (r *HtsgetRequest) SetFormat(format string) {
	r.format = format
}

func (r *HtsgetRequest) GetFormat() string {
	return r.format
}

func (r *HtsgetRequest) SetClass(class string) {
	r.class = class
}

func (r *HtsgetRequest) GetClass() string {
	return r.class
}

func (r *HtsgetRequest) SetReferenceName(referenceName string) {
	r.referenceName = referenceName
}

func (r *HtsgetRequest) GetReferenceName() string {
	return r.referenceName
}

func (r *HtsgetRequest) SetStart(start int) {
	r.start = start
}

func (r *HtsgetRequest) GetStart() int {
	return r.start
}

func (r *HtsgetRequest) SetEnd(end int) {
	r.end = end
}

func (r *HtsgetRequest) GetEnd() int {
	return r.end
}

func (r *HtsgetRequest) SetFields(fields []string) {
	r.fields = fields
}

func (r *HtsgetRequest) GetFields() []string {
	return r.fields
}

func (r *HtsgetRequest) SetTags(tags []string) {
	r.tags = tags
}

func (r *HtsgetRequest) GetTags() []string {
	return r.tags
}

func (r *HtsgetRequest) SetNoTags(noTags []string) {
	r.noTags = noTags
}

func (r *HtsgetRequest) GetNoTags() []string {
	return r.noTags
}

func (r *HtsgetRequest) SetRegions(regions []*Region) {
	r.regions = regions
}

func (r *HtsgetRequest) AddRegion(region *Region) {
	r.regions = append(r.regions, region)
}

func (r *HtsgetRequest) GetRegions() []*Region {
	return r.regions
}

func (r *HtsgetRequest) SetHtsgetBlockClass(htsgetBlockClass string) {
	r.htsgetBlockClass = htsgetBlockClass
}

func (r *HtsgetRequest) GetHtsgetBlockClass() string {
	return r.htsgetBlockClass
}

func (r *HtsgetRequest) SetHtsgetCurrentBlock(htsgetCurrentBlock string) {
	r.htsgetCurrentBlock = htsgetCurrentBlock
}

func (r *HtsgetRequest) GetHtsgetCurrentBlock() string {
	return r.htsgetCurrentBlock
}

func (r *HtsgetRequest) SetHtsgetTotalBlocks(htsgetTotalBlocks string) {
	r.htsgetTotalBlocks = htsgetTotalBlocks
}

func (r *HtsgetRequest) GetHtsgetTotalBlocks() string {
	return r.htsgetTotalBlocks
}

func (r *HtsgetRequest) SetHtsgetFilePath(htsgetFilePath string) {
	r.htsgetFilePath = htsgetFilePath
}

func (r *HtsgetRequest) GetHtsgetFilePath() string {
	return r.htsgetFilePath
}

func (r *HtsgetRequest) SetHtsgetRange(htsgetRange string) {
	r.htsgetRange = htsgetRange
}

func (r *HtsgetRequest) GetHtsgetRange() string {
	return r.htsgetRange
}

/* OTHER API METHODS */

func (r *HtsgetRequest) isDefaultString(val string, def string) bool {
	return val == def
}

func (r *HtsgetRequest) isDefaultInt(val int, def int) bool {
	return val == def
}

// isDefaultList checks if a list parameter value matches the default,
// unspecfied value, thereby indicating that the parameter was not specified
// in the HTTP request
func (r *HtsgetRequest) isDefaultList(val []string, def []string) bool {
	if len(val) != len(def) {
		return false
	}
	for i := 0; i < len(val); i++ {
		if val[i] != def[i] {
			return false
		}
	}
	return true
}

// HeaderOnlyRequested checks if the client request is only for the header
func (r *HtsgetRequest) HeaderOnlyRequested() bool {
	return r.GetClass() == htsconstants.ClassHeader
}

// UnplacedUnmappedReadsRequested checks if the client request is for unplaced,
// unmapped reads
func (r *HtsgetRequest) UnplacedUnmappedReadsRequested() bool {
	return r.GetReferenceName() == "*"
}

func (r *HtsgetRequest) ReferenceNameRequested() bool {
	return !r.isDefaultString(r.GetReferenceName(), defaultReferenceName)
}

func (r *HtsgetRequest) StartRequested() bool {
	return !r.isDefaultInt(r.GetStart(), defaultStart)
}

func (r *HtsgetRequest) EndRequested() bool {
	return !r.isDefaultInt(r.GetEnd(), defaultEnd)
}

func (r *HtsgetRequest) NRegions() int {
	return len(r.GetRegions())
}

// AllRegionsRequested checks if the client request is for all chromosomal
// regions in the file (ie. referenceName not specified)
func (r *HtsgetRequest) AllRegionsRequested() bool {
	return r.NRegions() == 0
}

// AllFieldsRequested checks if all fields were requested by the client. all
// fields are requested if the client does not specify the 'fields' parameter
func (r *HtsgetRequest) AllFieldsRequested() bool {
	return r.isDefaultList(r.GetFields(), defaultFields)
}

// TagsNotSpecified checks if the tags parameter was not provided in the HTTP
// request
func (r *HtsgetRequest) TagsNotSpecified() bool {
	return r.isDefaultList(r.GetTags(), defaultTags)
}

// NoTagsNotSpecified checks if the notags parameter was not provided in the
// HTTP request
func (r *HtsgetRequest) NoTagsNotSpecified() bool {
	return r.isDefaultList(r.GetNoTags(), defaultNoTags)
}

// AllTagsRequested checks if all tags were requested by the client. all tags
// are requested if the request specifies neither the 'tags' or 'notags' parameters
func (r *HtsgetRequest) AllTagsRequested() bool {
	return r.TagsNotSpecified() && r.NoTagsNotSpecified()
}

func (r *HtsgetRequest) IsFinalBlock() bool {
	current, _ := strconv.Atoi(r.GetHtsgetCurrentBlock())
	total, _ := strconv.Atoi(r.GetHtsgetTotalBlocks())
	return current == total-1
}

// ConstructDataEndpointURL for a given htsget request object, return the url
// that will redirect the client to the correct data download endpoint with
// all necessary parameters and headers provided
func (r *HtsgetRequest) ConstructDataEndpointURL(useRegion bool, regionI int) (string, error) {
	host := htsconfig.GetHost()
	dataEndpointPath := r.GetEndpoint().DataEndpointPath()
	dataEndpoint, err := url.Parse(htsutils.RemoveTrailingSlash(host) + dataEndpointPath + r.GetID())
	if err != nil {
		return "", err
	}

	// add query params
	query := dataEndpoint.Query()
	if r.HeaderOnlyRequested() {
		query.Set("class", r.GetClass())
	}

	if useRegion {
		region := r.GetRegions()[regionI]
		if region.ReferenceNameRequested() {
			query.Set("referenceName", region.GetReferenceName())
		}
		if region.StartRequested() {
			query.Set("start", region.StartString())
		}
		if region.EndRequested() {
			query.Set("end", region.EndString())
		}
	}

	if !r.AllFieldsRequested() {
		f := strings.Join(r.GetFields(), ",")
		query.Set("fields", f)
	}
	if !r.TagsNotSpecified() {
		t := strings.Join(r.GetTags(), ",")
		query.Set("tags", t)
	}
	if !r.NoTagsNotSpecified() {
		nt := strings.Join(r.GetNoTags(), ",")
		query.Set("notags", nt)
	}
	dataEndpoint.RawQuery = query.Encode()
	return dataEndpoint.String(), nil
}

func (r *HtsgetRequest) GetDataSourceRegistry() *htsconfig.DataSourceRegistry {
	return htsconfig.GetDataSourceRegistry(r.GetEndpoint())
}

func (r *HtsgetRequest) GetServiceInfo() *htsconfig.ServiceInfo {
	return htsconfig.GetServiceInfo(r.GetEndpoint())
}
