// Package htsformats manipulates bioinformatic data encountered by htsget
//
// Module samrecord_test tests region samrecord
package htsformats

import (
	"strings"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

var samrecordEmitCustomFieldsTC = []struct {
	line          string
	emittedFields []string
	exp           string
}{
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		[]string{"QNAME", "FLAG", "RNAME"},
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t0\t255\t*\t*\t0\t0\t*\t*",
	},
}

var samrecordEmitCustomTagsTC = []struct {
	line, tags, notags, exp string
}{
	// specify neither tags nor notags
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"ALL",
		"NONE",
		"NH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
	},
	// specify tags, but not notags
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"HI",
		"NONE",
		"HI:i:1",
	},
	// specify notags, but not tags
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"ALL",
		"HI",
		"NH:i:2\tNM:i:0\tMD:Z:100",
	},
}

var samrecordCustomEmitTC = []struct {
	line, fields, tags, notags, exp string
}{
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"ALL",
		"ALL",
		"NONE",
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
	},
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"QNAME,FLAG,RNAME",
		"ALL",
		"NONE",
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t0\t255\t*\t*\t0\t0\t*\t*\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
	},
	{
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tNH:i:2\tHI:i:1\tNM:i:0\tMD:Z:100",
		"ALL",
		"HI,NM",
		"NONE",
		"A00111:67:H3M5YDMXX:1:2407:21558:16094\t99\tchr1\t24613323\t3\t100M\t=\t24613553\t330\tCAATAAGGAATGTTGATCCAATAATTACATGGAGTCCATGGAATCCAGTAGCCATGAAGAATGTAGAACCATAGATACCATCTGAAATGGAGAATGATGT\tFFFFFFFFFFFFFFFFFF8FFFFF8FFFFFFFFFFFFFFFFFFFFFFFFF-FFFFFFF-FFFFFFFFF--F-FFFFFFFFF-FFFFF-FFF-F-FFFFFF\tHI:i:1\tNM:i:0",
	},
}

func TestSamRecordEmitCustomFields(t *testing.T) {
	for _, tc := range samrecordEmitCustomFieldsTC {
		samrecord := NewSAMRecord(tc.line)
		emittedFieldsArray := samrecord.emitCustomFields(tc.emittedFields)
		emittedFieldsString := strings.Join(emittedFieldsArray, "\t")
		if emittedFieldsString != tc.exp {
			t.Errorf("Expected: %s, Actual: %s", tc.exp, emittedFieldsString)
		}
	}
}

func TestSamRecordEmitCustomTags(t *testing.T) {
	for _, tc := range samrecordEmitCustomTagsTC {
		htsreq := htsrequest.NewHtsgetRequest()
		htsreq.AddListParam("tags", strings.Split(tc.tags, ","))
		htsreq.AddListParam("notags", strings.Split(tc.notags, ","))

		samrecord := NewSAMRecord(tc.line)
		emittedTagsArray := samrecord.emitCustomTags(htsreq)
		emittedTagsString := strings.Join(emittedTagsArray, "\t")
		if emittedTagsString != tc.exp {
			t.Errorf("Expected: %s, Actual: %s", tc.exp, emittedTagsString)
		}
	}
}

func TestSamRecordCustomEmit(t *testing.T) {
	for _, tc := range samrecordCustomEmitTC {
		htsreq := htsrequest.NewHtsgetRequest()
		htsreq.AddListParam("fields", strings.Split(tc.fields, ","))
		htsreq.AddListParam("tags", strings.Split(tc.tags, ","))
		htsreq.AddListParam("notags", strings.Split(tc.notags, ","))

		samrecord := NewSAMRecord(tc.line)
		emittedRecord := samrecord.CustomEmit(htsreq)
		if emittedRecord != tc.exp {
			t.Errorf("Expected: %s, Actual: %s", tc.exp, emittedRecord)
		}
	}
}
