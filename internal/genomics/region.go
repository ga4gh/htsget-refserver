// Package genomics deals with the traversal of genomic intervals
//
// Module region.go contains operations for parsing genomic intervals from
// coordinates
package genomics

// Region defines a chromosomal interval
//
// Attributes
//	Name (string): contig/chromosome name
//	Start (string): 0-index, inclusive start base position
//	End (string): 0-index, exclusive end base position
type Region struct {
	Name, Start, End string
}

// String get a representation of a genomic Region as a string
//
// Type: Region
// Returns
//	(string): Region interval as string
func (r *Region) String() string {
	if r.Name == "" || r.Name == "*" || r.Start == "-1" {
		return r.Name
	}
	if r.End == "-1" {
		return r.Name + ":" + r.Start
	} else {
		return r.Name + ":" + r.Start + "-" + r.End
	}
}
