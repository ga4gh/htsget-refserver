// Package htsformats manipulates bioinformatic data encountered by htsget
//
// Module region contains genomic intervals
package htsformats

// Region defines a simple genomic interval: contig name, start, and end position
type Region struct {
	Name, Start, End string
}

// String gets a representation of a genomic region
func (r *Region) String() string {
	if r.Start == "-1" && r.End == "-1" {
		return r.Name
	}
	if r.Start != "-1" && r.End == "-1" {
		return r.Name + ":" + r.Start
	}
	if r.Start == "-1" && r.End != "-1" {
		return r.Name + ":" + "0-" + r.End
	}
	return r.Name + ":" + r.Start + "-" + r.End
}

// ExportSamtools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (r *Region) ExportSamtools() string {
	return r.String()
}

// ExportBcftools exports the region in a manner compatible to how region requests
// are specified on the samtools command-line
func (r *Region) ExportBcftools() string {
	if r.Start == "-1" && r.End == "-1" {
		return r.Name
	}
	if r.Start != "-1" && r.End == "-1" {
		return r.Name + ":" + r.Start + "-"
	}
	if r.Start == "-1" && r.End != "-1" {
		return r.Name + ":" + "0-" + r.End
	}
	return r.Name + ":" + r.Start + "-" + r.End
}
