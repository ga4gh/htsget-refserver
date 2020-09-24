// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module region contains genomic intervals
package htsrequest

import "strconv"

// Region defines a simple genomic interval: contig name, start, and end position
type Region struct {
	ReferenceName string `json:"referenceName"`
	Start         int    `json:"start"`
	End           int    `json:"end"`
}

/* CONSTRUCTOR */

func NewRegion() *Region {
	return new(Region)
}

/* SETTERS AND GETTERS */

func (region *Region) SetReferenceName(referenceName string) {
	region.ReferenceName = referenceName
}

func (region *Region) GetReferenceName() string {
	return region.ReferenceName
}

func (region *Region) SetStart(start int) {
	region.Start = start
}

func (region *Region) GetStart() int {
	return region.Start
}

func (region *Region) SetEnd(end int) {
	region.End = end
}

func (region *Region) GetEnd() int {
	return region.End
}

func (region *Region) StartString() string {
	return strconv.Itoa(region.Start)
}

func (region *Region) EndString() string {
	return strconv.Itoa(region.End)
}

/* API METHODS */

func (region *Region) ReferenceNameRequested() bool {
	return !(region.GetReferenceName() == "")
}

func (region *Region) StartRequested() bool {
	return !(region.GetStart() == -1)
}

func (region *Region) EndRequested() bool {
	return !(region.GetEnd() == -1)
}

// String gets a representation of a genomic region
func (region *Region) String() string {
	if region.Start == -1 && region.End == -1 {
		return region.ReferenceName
	}
	if region.Start != -1 && region.End == -1 {
		return region.ReferenceName + ":" + region.StartString()
	}
	if region.Start == -1 && region.End != -1 {
		return region.ReferenceName + ":" + "0-" + region.EndString()
	}
	return region.ReferenceName + ":" + region.StartString() + "-" + region.EndString()
}

// ExportSamtools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (region *Region) ExportSamtools() string {
	return region.String()
}

// ExportBcftools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (region *Region) ExportBcftools() string {
	if region.Start == -1 && region.End == -1 {
		return region.ReferenceName
	}
	if region.Start != -1 && region.End == -1 {
		return region.ReferenceName + ":" + region.StartString() + "-"
	}
	if region.Start == -1 && region.End != -1 {
		return region.ReferenceName + ":" + "0-" + region.EndString()
	}
	return region.ReferenceName + ":" + region.StartString() + "-" + region.EndString()
}
