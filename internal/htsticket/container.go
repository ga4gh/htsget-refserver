// Package htsticket ...
package htsticket

// Container holds the file format, urls of files for the client,
// and optionally an MD5 digest resulting from the concatenation of url data blocks
type Container struct {
	Format string `json:"format"`
	URLS   []*URL `json:"urls"`
	MD5    string `json:"md5,omitempty"`
}

func NewContainer() *Container {
	return new(Container)
}

func (container *Container) setFormat(format string) *Container {
	container.Format = format
	return container
}

func (container *Container) SetFormatBam() *Container {
	container.setFormat("BAM")
	return container
}

func (container *Container) SetFormatVcf() *Container {
	container.setFormat("VCF")
	return container
}

func (container *Container) SetURLS(urls []*URL) *Container {
	container.URLS = urls
	return container
}
