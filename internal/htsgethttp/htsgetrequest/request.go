// Package htsgetrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
package htsgetrequest

// HtsgetRequest contains htsget-related parameters within two maps,
// one contains scalar values, the other contains list values
type HtsgetRequest struct {
	ScalarParams map[string]string
	ListParams   map[string][]string
}

// NewHtsgetRequest constructs a new HtsgetRequest struct instance
func NewHtsgetRequest() *HtsgetRequest {
	htsgetReq := new(HtsgetRequest)
	htsgetReq.ScalarParams = make(map[string]string)
	htsgetReq.ListParams = make(map[string][]string)
	return htsgetReq
}

// AddScalarParam adds a key-value pair to the scalar parameter map
func (htsgetReq *HtsgetRequest) AddScalarParam(key string, value string) {
	htsgetReq.ScalarParams[key] = value
}

// AddListParam adds a key-value pair to the list parameter map
func (htsgetReq *HtsgetRequest) AddListParam(key string, value []string) {
	htsgetReq.ListParams[key] = value
}

// Get retrieves a value from the scalar parameter map by its key
func (htsgetReq *HtsgetRequest) Get(key string) string {
	return htsgetReq.ScalarParams[key]
}

// GetList retrieves a value from the list parameter map by its key
func (htsgetReq *HtsgetRequest) GetList(key string) []string {
	return htsgetReq.ListParams[key]
}
