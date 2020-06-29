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
func (htsgetReq *HtsgetRequest) get(key string) string {
	return htsgetReq.ScalarParams[key]
}

// GetList retrieves a value from the list parameter map by its key
func (htsgetReq *HtsgetRequest) getList(key string) []string {
	return htsgetReq.ListParams[key]
}

// ID get request object id
// Type: *HtsgetRequest - HtsgetRequest struct
// Arguments: None
// Returns: string - request object id
func (htsgetReq *HtsgetRequest) ID() string {
	return htsgetReq.get("id")
}

func (htsgetReq *HtsgetRequest) Format() string {
	return htsgetReq.get("format")
}

func (htsgetReq *HtsgetRequest) Class() string {
	return htsgetReq.get("class")
}

func (htsgetReq *HtsgetRequest) ReferenceName() string {
	return htsgetReq.get("referenceName")
}

func (htsgetReq *HtsgetRequest) Start() string {
	return htsgetReq.get("start")
}

func (htsgetReq *HtsgetRequest) End() string {
	return htsgetReq.get("end")
}

func (htsgetReq *HtsgetRequest) HtsgetBlockClass() string {
	return htsgetReq.get("HtsgetBlockClass")
}

func (htsgetReq *HtsgetRequest) HtsgetBlockId() string {
	return htsgetReq.get("HtsgetBlockId")
}

func (htsgetReq *HtsgetRequest) HtsgetNumBlocks() string {
	return htsgetReq.get("HtsgetNumBlocks")
}

func (htsgetReq *HtsgetRequest) Fields() []string {
	return htsgetReq.getList("fields")
}

func (htsgetReq *HtsgetRequest) Tags() []string {
	return htsgetReq.getList("tags")
}

func (htsgetReq *HtsgetRequest) NoTags() []string {
	return htsgetReq.getList("notags")
}

func (htsgetReq *HtsgetRequest) isDefaultScalar(key string) bool {
	return htsgetReq.get(key) == defaultScalarParameterValues[key]
}

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

func (htsgetReq *HtsgetRequest) AllFieldsRequested() bool {
	return htsgetReq.isDefaultList("fields")
}

func (htsgetReq *HtsgetRequest) TagsNotSpecified() bool {
	return htsgetReq.isDefaultList("tags")
}

func (htsgetReq *HtsgetRequest) NoTagsNotSpecified() bool {
	return htsgetReq.isDefaultList("notags")
}

func (htsgetReq *HtsgetRequest) AllTagsRequested() bool {
	return htsgetReq.TagsNotSpecified() && htsgetReq.NoTagsNotSpecified()
}
