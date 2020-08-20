// Package htsgetformats contains operations for reading and manipulating data
// in bioinformatics file formats encountered by htsget
//
// Module samrecord.go contains operations for working with individual records
// (data lines) within a SAM/BAM/CRAM file. Each record represents a single
// read
package htsgetformats

import (
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsgetconfig/htsgetconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsgetutils"
)

// SAMRecord contains all columnar data found in a single SAM/BAM file record
//
// Attributes
//	line (string): the record string as it appears in the SAM/BAM file
//	columns ([]string): record data split according to column in the expected order
type SAMRecord struct {
	line    string
	columns []string
}

// NewSAMRecord instantiates a new SAMRecord instance
//
// Arguments
//	line (string): the original sam record line
// Returns
//	(*SAMRecord): new SAMRecord instance
func NewSAMRecord(line string) *SAMRecord {
	samRecord := new(SAMRecord)
	samRecord.line = line
	samRecord.columns = strings.Split(samRecord.line, "\t")
	return samRecord
}

// emitCustomFields includes/excludes fields based on the 'fields' parameter
// of the HTTP request. excluded fields are replaced with the appropriate
// exclusion/missing data values
//
// Type: SAMRecord
// Arguments
//	fields ([]string): requested fields to include
// Returns
//	([]string): new SAM columns, with values included/excluded based on request
func (samRecord *SAMRecord) emitCustomFields(fields []string) []string {

	n := htsgetconstants.BamFieldsN
	emittedFields := make([]string, n)
	toEmitByField := make([]bool, n)

	// initializes an array, saying that each column will be excluded
	for i := 0; i < n; i++ {
		toEmitByField[i] = false
	}

	// set each column in the fields list to 'true,' ie. the true value will
	// be included
	for i := 0; i < len(fields); i++ {
		toEmitByField[htsgetconstants.BamFields[fields[i]]] = true
	}

	// if the value for a column is true, add the true value to the list,
	// otherwise add the correct excluded value for that column
	for i := 0; i < htsgetconstants.BamFieldsN; i++ {
		if toEmitByField[i] {
			emittedFields[i] = samRecord.columns[i]
		} else {
			emittedFields[i] = htsgetconstants.BamExcludedValues[i]
		}
	}

	return emittedFields
}

// emitCustomTags includes/excludes tags based on the 'tags' and 'notags'
// parameters of the HTTP request. exluded tags are removed from the output
//
// Type: SAMRecord
// Arguments
//	htsgetReq (*HtsgetReqest): the htsget request object
// Returns
//	([]string): new SAM tag columns, only including those requested
func (samRecord *SAMRecord) emitCustomTags(htsgetReq *htsgetrequest.HtsgetRequest) []string {

	tags := htsgetReq.Tags()
	notags := htsgetReq.NoTags()
	n := htsgetconstants.BamFieldsN
	emittedTags := make([]string, 0)

	// for each tag column:
	for i := n; i < len(samRecord.columns); i++ {
		tag := samRecord.columns[i]
		tagName := htsgetutils.GetTagName(tag)

		// if the 'tags' parameter is not specified, then the client has requested
		// all tags minus those in notags. if it is specifed, the client only
		// wants the tags specified by 'tags'
		toEmit := false
		if htsgetReq.TagsNotSpecified() {
			toEmit = true
		}

		if htsgetutils.IsItemInArray(tagName, tags) {
			toEmit = true
		}

		if htsgetutils.IsItemInArray(tagName, notags) {
			toEmit = false
		}

		if toEmit {
			emittedTags = append(emittedTags, tag)
		}
	}

	return emittedTags
}

// CustomEmit emits a new SAM record from an existing record, with only the
// requested fields and tags included
//
// Type: SAMRecord
// Arguments
//	htsgetReq (*HtsgetRequest): the htsget request object
// Returns
//	(string): the custom representation of the record based on requested fields and tags
func (samRecord *SAMRecord) CustomEmit(htsgetReq *htsgetrequest.HtsgetRequest) string {

	var emittedFields []string
	var emittedTags []string

	// only run the custom emit fields function if the client has requested
	// specific fields
	if htsgetReq.AllFieldsRequested() {
		emittedFields = samRecord.columns[0:11]
	} else {
		emittedFields = samRecord.emitCustomFields(htsgetReq.Fields())
	}

	// only run the custom emit tags function if the client has requested
	// specific tags be included or excluded
	if htsgetReq.AllTagsRequested() {
		emittedTags = samRecord.columns[11:]
	} else {
		emittedTags = samRecord.emitCustomTags(htsgetReq)
	}

	for i := 0; i < len(emittedTags); i++ {
		emittedFields = append(emittedFields, emittedTags[i])
	}

	return strings.Join(emittedFields, "\t")
}
