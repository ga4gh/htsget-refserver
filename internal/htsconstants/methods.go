package htsconstants

type HttpMethod int

const (
	GetMethod  HttpMethod = 0
	PostMethod HttpMethod = 1
)

const (
	GetMethodS  string = "GET"
	PostMethodS string = "POST"
)

var httpMethodStringMap = map[HttpMethod]string{
	GetMethod:  GetMethodS,
	PostMethod: PostMethodS,
}

func (e HttpMethod) String() string {
	return httpMethodStringMap[e]
}
