// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module request defines structs and operations for a mature htsget
// request, which holds all htsget-related parameters and has insight into
// what information was requested by the client
package htsrequest

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"

	log "github.com/sirupsen/logrus"
)

// HtsgetRequest contains htsget-related parameters
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
	headers            http.Header
}

// NewHtsgetRequest instantiates a new HtsgetRequest instance
func NewHtsgetRequest() *HtsgetRequest {
	r := new(HtsgetRequest)
	r.SetRegions([]*Region{})
	return r
}

// SetEndpoint sets the API endpoint associated with request
func (r *HtsgetRequest) SetEndpoint(endpoint htsconstants.APIEndpoint) {
	r.endpoint = endpoint
}

// GetEndpoint retrieves the API endpoint associated with request
func (r *HtsgetRequest) GetEndpoint() htsconstants.APIEndpoint {
	return r.endpoint
}

// SetID sets request ID
func (r *HtsgetRequest) SetID(id string) {
	r.id = id
}

// GetID retrieves request ID
func (r *HtsgetRequest) GetID() string {
	return r.id
}

// SetFormat sets the requested file format
func (r *HtsgetRequest) SetFormat(format string) {
	r.format = format
}

// GetFormat retrieves the requested file format
func (r *HtsgetRequest) GetFormat() string {
	return r.format
}

// SetClass sets the requested class (ie. for header requests)
func (r *HtsgetRequest) SetClass(class string) {
	r.class = class
}

// GetClass retrieves the requested class
func (r *HtsgetRequest) GetClass() string {
	return r.class
}

// SetReferenceName sets the requested chromosome/reference sequence name
func (r *HtsgetRequest) SetReferenceName(referenceName string) {
	r.referenceName = referenceName
}

// GetReferenceName retrieves the requested chromosome/reference sequence name
func (r *HtsgetRequest) GetReferenceName() string {
	return r.referenceName
}

// SetStart sets the requested region start position
func (r *HtsgetRequest) SetStart(start int) {
	r.start = start
}

// GetStart retrieves the requested region start position
func (r *HtsgetRequest) GetStart() int {
	return r.start
}

// SetEnd sets the requested region end position
func (r *HtsgetRequest) SetEnd(end int) {
	r.end = end
}

// GetEnd retrieves the requested region end position
func (r *HtsgetRequest) GetEnd() int {
	return r.end
}

// SetFields sets the requested emitted fields
func (r *HtsgetRequest) SetFields(fields []string) {
	r.fields = fields
}

// GetFields retrieves the requested emitted fields
func (r *HtsgetRequest) GetFields() []string {
	return r.fields
}

// SetTags sets the requested emitted tags
func (r *HtsgetRequest) SetTags(tags []string) {
	r.tags = tags
}

// GetTags retrieves the requested emitted tags
func (r *HtsgetRequest) GetTags() []string {
	return r.tags
}

// SetNoTags sets the requested set of tags to be excluded from results
func (r *HtsgetRequest) SetNoTags(noTags []string) {
	r.noTags = noTags
}

// GetNoTags retrieves the requested set of tags to be excluded from results
func (r *HtsgetRequest) GetNoTags() []string {
	return r.noTags
}

// SetRegions sets the requested list of genomic regions to be returned
func (r *HtsgetRequest) SetRegions(regions []*Region) {
	r.regions = regions
}

// AddRegion adds a single region to the list of requested genomic regions
func (r *HtsgetRequest) AddRegion(region *Region) {
	r.regions = append(r.regions, region)
}

// GetRegions retrieves the requested list of genomic regions
func (r *HtsgetRequest) GetRegions() []*Region {
	return r.regions
}

// SetHtsgetBlockClass sets the request block class
func (r *HtsgetRequest) SetHtsgetBlockClass(htsgetBlockClass string) {
	r.htsgetBlockClass = htsgetBlockClass
}

// GetHtsgetBlockClass retrieves the requested block class
func (r *HtsgetRequest) GetHtsgetBlockClass() string {
	return r.htsgetBlockClass
}

// SetHtsgetCurrentBlock sets the current url data block number
func (r *HtsgetRequest) SetHtsgetCurrentBlock(htsgetCurrentBlock string) {
	r.htsgetCurrentBlock = htsgetCurrentBlock
}

// GetHtsgetCurrentBlock retrieves the current url data block number
func (r *HtsgetRequest) GetHtsgetCurrentBlock() string {
	return r.htsgetCurrentBlock
}

// SetHtsgetTotalBlocks sets the total number of expected url data blocks /
// fileparts for a single file download
func (r *HtsgetRequest) SetHtsgetTotalBlocks(htsgetTotalBlocks string) {
	r.htsgetTotalBlocks = htsgetTotalBlocks
}

// GetHtsgetTotalBlocks retrieves the total number of expected url data blocks /
// fileparts for a single file download
func (r *HtsgetRequest) GetHtsgetTotalBlocks() string {
	return r.htsgetTotalBlocks
}

// SetHtsgetFilePath sets the path to the requested source file
func (r *HtsgetRequest) SetHtsgetFilePath(htsgetFilePath string) {
	r.htsgetFilePath = htsgetFilePath
}

// GetHtsgetFilePath retrieves the path to the requested source file
func (r *HtsgetRequest) GetHtsgetFilePath() string {
	return r.htsgetFilePath
}

// SetHtsgetRange sets the byte range for a requested file
func (r *HtsgetRequest) SetHtsgetRange(htsgetRange string) {
	r.htsgetRange = htsgetRange
}

// GetHtsgetRange retrieves the byte range for a requested file
func (r *HtsgetRequest) GetHtsgetRange() string {
	return r.htsgetRange
}

// isDefaultString checks if a given string property matches the expected default value
func (r *HtsgetRequest) isDefaultString(val string, def string) bool {
	return val == def
}

// isDefaultInt checks if a given int property matches the expected default value
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

// ReferenceNameRequested checks whether a reference name was specified in the
// request
func (r *HtsgetRequest) ReferenceNameRequested() bool {
	return !r.isDefaultString(r.GetReferenceName(), defaultReferenceName)
}

// StartRequested checks whether a genomic start position was specified in the request
func (r *HtsgetRequest) StartRequested() bool {
	return !r.isDefaultInt(r.GetStart(), defaultStart)
}

// EndRequested checks whether a genomic end position was specified in the request
func (r *HtsgetRequest) EndRequested() bool {
	return !r.isDefaultInt(r.GetEnd(), defaultEnd)
}

// NRegions returns the number of requested genomic loci
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

// IsHeaderBlock checks whether the current data block / filepart represents
// the header of the genomic file
func (r *HtsgetRequest) IsHeaderBlock() bool {
	current, _ := strconv.Atoi(r.GetHtsgetCurrentBlock())
	return current == 0
}

// IsFinalBlock checks whether the current data block / filepart is the final
// block in a list of blocks, altogether constituting a single file
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
		log.Errorf("error parsing in ConstructDataEndpointURL, %v", err)
		return "", err
	}

	// add query params
	query := dataEndpoint.Query()
	if r.HeaderOnlyRequested() {
		query.Set("class", r.GetClass())
	}

	useRegion = false // temporary solution: return full file when doing partial request

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

// GetDataSourceRegistry retrieves the data sources associated with the endpoint
func (r *HtsgetRequest) GetDataSourceRegistry() *htsconfig.DataSourceRegistry {
	return htsconfig.GetDataSourceRegistry(r.GetEndpoint())
}

// GetServiceInfo retrieves the service info object associated with the endpoint
func (r *HtsgetRequest) GetServiceInfo() *htsconfig.ServiceInfo {
	return htsconfig.GetServiceInfo(r.GetEndpoint())
}
