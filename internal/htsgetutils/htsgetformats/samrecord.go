package htsgetformats

import (
	"fmt"
	"strings"

	"github.com/ga4gh/htsget-refserver/internal/config"

	"github.com/ga4gh/htsget-refserver/internal/htsgethttp/htsgetrequest"
)

type SAMRecord struct {
	line    string
	columns []string
}

func NewSAMRecord(line string) *SAMRecord {
	samRecord := new(SAMRecord)
	samRecord.line = line
	samRecord.columns = strings.Split(samRecord.line, "\t")
	return samRecord
}

func (samRecord *SAMRecord) emitCustomFields(fields []string) []string {

	n := config.N_BAM_FIELDS
	emittedFields := make([]string, config.N_BAM_FIELDS)
	toEmitByField := make([]bool, config.N_BAM_FIELDS)

	for i := 0; i < n; i++ {
		toEmitByField[i] = false
	}

	for i := 0; i < len(fields); i++ {
		toEmitByField[config.BAM_FIELDS[fields[i]]] = true
	}

	for i := 0; i < config.N_BAM_FIELDS; i++ {
		if toEmitByField[i] {
			emittedFields[i] = samRecord.columns[i]
		} else {
			emittedFields[i] = config.BAM_EXCLUDED_VALUES[i]
		}
	}

	return emittedFields
}

func (samRecord *SAMRecord) emitCustomTags(tags []string, notags []string) []string {
	return []string{
		"TAG",
	}
}

func (samRecord *SAMRecord) CustomEmit(htsgetReq *htsgetrequest.HtsgetRequest) string {

	var emittedFields []string
	var emittedTags []string

	if htsgetReq.AllFieldsRequested() {
		emittedFields = samRecord.columns[0:11]
	} else {
		emittedFields = samRecord.emitCustomFields(htsgetReq.Fields())
	}

	fmt.Println(emittedFields)

	if htsgetReq.AllTagsRequested() {
		emittedTags = samRecord.columns[11:]
	} else {
		emittedTags = samRecord.emitCustomTags(htsgetReq.Tags(), htsgetReq.NoTags())
	}

	for i := 0; i < len(emittedTags); i++ {
		emittedFields = append(emittedFields, emittedTags[i])
	}

	return strings.Join(emittedFields, "\t")
}
