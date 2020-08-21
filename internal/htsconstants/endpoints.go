package htsconstants

type ServerEndpoint int

const (
	ReadsTicket         ServerEndpoint = 0
	ReadsData           ServerEndpoint = 1
	ReadsServiceInfo    ServerEndpoint = 2
	VariantsTicket      ServerEndpoint = 3
	VariantsData        ServerEndpoint = 4
	VariantsServiceInfo ServerEndpoint = 5
	FileBytes           ServerEndpoint = 6
)

const (
	ReadsTicketS         string = "/reads/{id}"
	ReadsDataS           string = "/reads/data/{id}"
	ReadsServiceInfoS    string = "/reads/service-info"
	VariantsTicketS      string = "/variants/{id}"
	VariantsDataS        string = "/variants/data/{id}"
	VariantsServiceInfoS string = "/variants/service-info"
	FileBytesS           string = "/file-bytes"
)

var htsEndpointStringMap = map[ServerEndpoint]string{
	ReadsTicket:         ReadsTicketS,
	ReadsData:           ReadsDataS,
	ReadsServiceInfo:    ReadsServiceInfoS,
	VariantsTicket:      VariantsTicketS,
	VariantsData:        VariantsDataS,
	VariantsServiceInfo: VariantsServiceInfoS,
	FileBytes:           FileBytesS,
}

func (e ServerEndpoint) String() string {
	return htsEndpointStringMap[e]
}
