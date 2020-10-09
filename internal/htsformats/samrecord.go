// Package htsformats manipulates bioinformatic data encountered by htsget
//
// Module samrecord.go contains operations for working with individual records
// (data lines) within a SAM/BAM/CRAM file. Each record represents a single
// read
package htsformats

import (
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/ga4gh/htsget-refserver/internal/htsutils"
)

// SAMRecord contains all columnar data found in a single SAM/BAM file record,
// as both the raw string (line) and separated into individual columns (columns)
type SAMRecord struct {
	line    string
	columns []string
}

// NewSAMRecord instantiates and returns a SAMRecord object from a single line
// of a SAM/BAM file
func NewSAMRecord(line string) *SAMRecord {
	samRecord := new(SAMRecord)
	samRecord.line = line
	samRecord.columns = strings.Split(samRecord.line, "\t")
	return samRecord
}

// emitCustomFields includes/excludes fields based on the custom inclusion
// criteria. excluded fields are replaced with the appropriate exclusion/missing
// data values. The fields list indicates what fields to emit
func (samRecord *SAMRecord) emitCustomFields(fields []string) []string {

	n := htsconstants.BamFieldsN
	emittedFields := make([]string, n)
	toEmitByField := make([]bool, n)

	// initializes an array, saying that each column will be excluded
	for i := 0; i < n; i++ {
		toEmitByField[i] = false
	}

	// set each column in the fields list to 'true,' ie. the true value will
	// be included
	for i := 0; i < len(fields); i++ {
		toEmitByField[htsconstants.BamFields[fields[i]]] = true
	}

	// if the value for a column is true, add the true value to the list,
	// otherwise add the correct excluded value for that column
	for i := 0; i < htsconstants.BamFieldsN; i++ {
		if toEmitByField[i] {
			emittedFields[i] = samRecord.columns[i]
		} else {
			emittedFields[i] = htsconstants.BamExcludedValues[i]
		}
	}

	return emittedFields
}

// emitCustomTags includes/excludes tags based on custom criteria (ie. tags and
// notags parameters of HTTP request). exluded tags are removed from the output
func (samRecord *SAMRecord) emitCustomTags(htsgetReq *htsrequest.HtsgetRequest) []string {

	tags := htsgetReq.GetTags()
	notags := htsgetReq.GetNoTags()
	n := htsconstants.BamFieldsN
	emittedTags := make([]string, 0)

	// for each tag column:
	for i := n; i < len(samRecord.columns); i++ {
		tag := samRecord.columns[i]
		tagName := htsutils.GetTagName(tag)

		// if the 'tags' parameter is not specified, then the client has requested
		// all tags minus those in notags. if it is specifed, the client only
		// wants the tags specified by 'tags'
		toEmit := false
		if htsgetReq.TagsNotSpecified() {
			toEmit = true
		}

		if htsutils.IsItemInArray(tagName, tags) {
			toEmit = true
		}

		if htsutils.IsItemInArray(tagName, notags) {
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
func (samRecord *SAMRecord) CustomEmit(htsgetReq *htsrequest.HtsgetRequest) string {

	var emittedFields []string
	var emittedTags []string

	// only run the custom emit fields function if the client has requested
	// specific fields
	if htsgetReq.AllFieldsRequested() {
		emittedFields = samRecord.columns[0:11]
	} else {
		emittedFields = samRecord.emitCustomFields(htsgetReq.GetFields())
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
