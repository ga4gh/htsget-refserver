// Package htsconstants contains program constants
//
// Module methods contains constants relating to HTTP/REST methods
package htsconstants

// HTTPMethod enum for http methods
type HTTPMethod int

// enum values for HTTPMethod
const (
	GetMethod  HTTPMethod = 0
	PostMethod HTTPMethod = 1
)

// string representations of HTTPMethod enum
const (
	GetMethodS  string = "GET"
	PostMethodS string = "POST"
)

// httpMethodStringMap maps HTTPMethod
var httpMethodStringMap = map[HTTPMethod]string{
	GetMethod:  GetMethodS,
	PostMethod: PostMethodS,
}

// String gets the string representation of an HTTPMethod enum instance
func (e HTTPMethod) String() string {
	return httpMethodStringMap[e]
}
