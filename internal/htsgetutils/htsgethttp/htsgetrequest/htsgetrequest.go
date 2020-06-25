package htsgetrequest

type HtsgetRequest struct {
	Scalars map[string]string
	Lists   map[string][]string
}

func New() *HtsgetRequest {
	htsgetReq := new(HtsgetRequest)
	htsgetReq.Scalars = make(map[string]string)
	htsgetReq.Lists = make(map[string][]string)
	return htsgetReq
}

func (htsgetReq *HtsgetRequest) AddToScalars(key string, value string) {
	htsgetReq.Scalars[key] = value
}

func (htsgetReq *HtsgetRequest) AddToLists(key string, value []string) {
	htsgetReq.Lists[key] = value
}

func (htsgetReq *HtsgetRequest) GetScalar(key string) string {
	return htsgetReq.Scalars[key]
}

func (htsgetReq *HtsgetRequest) GetList(key string) []string {
	return htsgetReq.Lists[key]
}
