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

func (r *Region) StartString() string {
	return strconv.Itoa(r.Start)
}

func (r *Region) EndString() string {
	return strconv.Itoa(r.End)
}

// String gets a representation of a genomic region
func (r *Region) String() string {
	if r.Start == -1 && r.End == -1 {
		return r.ReferenceName
	}
	if r.Start != -1 && r.End == -1 {
		return r.ReferenceName + ":" + r.StartString()
	}
	if r.Start == -1 && r.End != -1 {
		return r.ReferenceName + ":" + "0-" + r.EndString()
	}
	return r.ReferenceName + ":" + r.StartString() + "-" + r.EndString()
}

// ExportSamtools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (r *Region) ExportSamtools() string {
	return r.String()
}

// ExportBcftools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (r *Region) ExportBcftools() string {
	if r.Start == -1 && r.End == -1 {
		return r.ReferenceName
	}
	if r.Start != -1 && r.End == -1 {
		return r.ReferenceName + ":" + r.StartString() + "-"
	}
	if r.Start == -1 && r.End != -1 {
		return r.ReferenceName + ":" + "0-" + r.EndString()
	}
	return r.ReferenceName + ":" + r.StartString() + "-" + r.EndString()
}
