// Package htsconstants contains program constants
//
// Module endpoints contains constants relating to htsget-specific endpoints
package htsconstants

// ServerEndpoint enum for different htsget-specific API routes
type ServerEndpoint int

// enum values for ServerEndpoint
const (
	ReadsTicket         ServerEndpoint = 0
	ReadsData           ServerEndpoint = 1
	ReadsServiceInfo    ServerEndpoint = 2
	VariantsTicket      ServerEndpoint = 3
	VariantsData        ServerEndpoint = 4
	VariantsServiceInfo ServerEndpoint = 5
	FileBytes           ServerEndpoint = 6
)

// string representations of ServerEndpoint enum
const (
	ReadsTicketS         string = "/reads/{id}"
	ReadsDataS           string = "/reads/data/{id}"
	ReadsServiceInfoS    string = "/reads/service-info"
	VariantsTicketS      string = "/variants/{id}"
	VariantsDataS        string = "/variants/data/{id}"
	VariantsServiceInfoS string = "/variants/service-info"
	FileBytesS           string = "/file-bytes"
)

// htsEndpointStringMap maps enum int values to string representation
var htsEndpointStringMap = map[ServerEndpoint]string{
	ReadsTicket:         ReadsTicketS,
	ReadsData:           ReadsDataS,
	ReadsServiceInfo:    ReadsServiceInfoS,
	VariantsTicket:      VariantsTicketS,
	VariantsData:        VariantsDataS,
	VariantsServiceInfo: VariantsServiceInfoS,
	FileBytes:           FileBytesS,
}

// String gets the string representation of a ServerEndpoint enum value
func (e ServerEndpoint) String() string {
	return htsEndpointStringMap[e]
}
