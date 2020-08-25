// Package htsconstants contains program constants
//
// Module endpoints contains constants relating to htsget-specific endpoints
package htsconstants

// APIEndpoint enum for different htsget-specific API routes
type APIEndpoint int

// enum values for ServerEndpoint
const (
	APIEndpointReadsTicket         APIEndpoint = 0
	APIEndpointReadsData           APIEndpoint = 1
	APIEndpointReadsServiceInfo    APIEndpoint = 2
	APIEndpointVariantsTicket      APIEndpoint = 3
	APIEndpointVariantsData        APIEndpoint = 4
	APIEndpointVariantsServiceInfo APIEndpoint = 5
	APIEndpointFileBytes           APIEndpoint = 6
)

// maps enum int values to string representation
var htsEndpointStringMap = map[APIEndpoint]string{
	APIEndpointReadsTicket:         "/reads/{id}",
	APIEndpointReadsData:           "/reads/data/{id}",
	APIEndpointReadsServiceInfo:    "/reads/service-info",
	APIEndpointVariantsTicket:      "/variants/{id}",
	APIEndpointVariantsData:        "/variants/data/{id}",
	APIEndpointVariantsServiceInfo: "/variants/service-info",
	APIEndpointFileBytes:           "/file-bytes",
}

// maps ticket endpoints to their corresponding data endpoint prefixes
var ticketEndpointToDataEndpointPathMap = map[APIEndpoint]string{
	APIEndpointReadsTicket:    "/reads/data/",
	APIEndpointVariantsTicket: "/variants/data/",
}

// maps endpoints to allowed format values
var endpointToEnabledFormatsMap = map[APIEndpoint][]string{
	APIEndpointReadsTicket:    []string{FormatBam /*, FormatCram */},
	APIEndpointReadsData:      []string{FormatBam /*, FormatCram */},
	APIEndpointVariantsTicket: []string{FormatVcf /*, FormatBcf */},
	APIEndpointVariantsData:   []string{FormatVcf /*, FormatBcf */},
}

// String gets the string representation of a ServerEndpoint enum value
func (e APIEndpoint) String() string {
	return htsEndpointStringMap[e]
}

// DataEndpointPath gets the corresponding data endpoint prefix for a given
// ticket APIEndpoint
func (e APIEndpoint) DataEndpointPath() string {
	return ticketEndpointToDataEndpointPathMap[e]
}

// AllowedFormats gets the acceptable requested formats based on the API Endpoint
func (e APIEndpoint) AllowedFormats() []string {
	return endpointToEnabledFormatsMap[e]
}
