// Package htsticket produces the htsget JSON response ticket
//
// Module container holds all attributes beneath the base ticket object
package htsticket

// Container holds the file format, urls of files for the client, and optionally
// an MD5 digest resulting from the concatenation of url data blocks
type Container struct {
	Format string `json:"format"`
	URLS   []*URL `json:"urls"`
	MD5    string `json:"md5,omitempty"`
}

// NewContainer instantiates and returns an empty ticket container
func NewContainer() *Container {
	return new(Container)
}

// setFormat sets the format of the container
func (container *Container) setFormat(format string) *Container {
	container.Format = format
	return container
}

// SetFormatBam sets the container format to BAM
func (container *Container) SetFormatBam() *Container {
	container.setFormat("BAM")
	return container
}

// SetFormatCram sets the container format to CRAM
func (container *Container) SetFormatCram() *Container {
	container.setFormat("CRAM")
	return container
}

// SetFormatVcf sets the container format to VCF
func (container *Container) SetFormatVcf() *Container {
	container.setFormat("VCF")
	return container
}

// SetFormatBcf sets the container format to BCF
func (container *Container) SetFormatBcf() *Container {
	container.setFormat("BCF")
	return container
}

// SetURLS sets the container's data download urls
func (container *Container) SetURLS(urls []*URL) *Container {
	container.URLS = urls
	return container
}
