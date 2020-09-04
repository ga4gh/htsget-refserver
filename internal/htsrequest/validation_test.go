package htsrequest

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"

	"github.com/stretchr/testify/assert"
)

var validateIDTC = []struct {
	endpoint htsconstants.APIEndpoint
	id       string
	exp      bool
}{
	{htsconstants.APIEndpointReadsTicket, "NoDataSource.00001", false},
	{htsconstants.APIEndpointReadsTicket, "tabulamuris.NoID", false},
}

var validateFormatTC = []struct {
	endpoint htsconstants.APIEndpoint
	format   string
	exp      bool
}{
	{htsconstants.APIEndpointReadsTicket, "BAM", true},
	{htsconstants.APIEndpointReadsTicket, "CRAM", false},
	{htsconstants.APIEndpointVariantsTicket, "VCF", true},
	{htsconstants.APIEndpointVariantsTicket, "BAM", false},
}

var validateClassTC = []struct {
	class string
	exp   bool
}{
	{"header", true},
	{"body", false},
	{"otherclass", false},
}

var validateReferenceNameTC = []struct {
	endpoint                 htsconstants.APIEndpoint
	id, class, referenceName string
	exp                      bool
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"header",
		"",
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"",
		"*",
		true,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"",
		"chr30",
		false,
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		"",
		"chr1",
		true,
	},
}

var validateStartTC = []struct {
	class, referenceName, start string
	exp                         bool
}{
	{"header", "", "10000", false},
	{"", "*", "10000", false},
	{"", "", "10000", false},
	{"", "chr1", "abc", false},
	{"", "chr1", "-10000", false},
	{"", "chr1", "10000", true},
	{"", "chr22", "220000", true},
	{"", "chrMT", "3000", true},
}

var validateEndTC = []struct {
	class, referenceName, start, end string
	exp                              bool
}{
	{"header", "", "10000", "20000", false},
	{"", "*", "10000", "20000", false},
	{"", "", "10000", "20000", false},
	{"", "chr1", "0", "abc", false},
	{"", "chr1", "0", "-10000", false},
	{"", "chr1", "100000", "50000", false},
	{"", "chr1", "100.50", "500.50", false},
	{"", "chr1", "100", "500", true},
}

var validateFieldsTC = []struct {
	class, fields string
	exp           bool
}{
	{"header", "", false},
	{"", "FOO,FLAG,QNAME", false},
	{"", "FLAG", true},
	{"", "TLEN,SEQ,QUAL,FLAG", true},
}

var validateTagsTC = []struct {
	class, tags string
	exp         bool
}{
	{"header", "NM,MD", false},
	{"", "NM,MD", true},
	{"", "HZ,MD,NM,HI", true},
}

var validateNoTagsTC = []struct {
	class  string
	tags   []string
	notags string
	exp    bool
}{
	{"header", []string{"NM", "MD"}, "", false},
	{"", []string{"NM", "MD"}, "", true},
	{"", []string{"NM", "MD"}, "MD", false},
}

func TestValidateID(t *testing.T) {
	for _, tc := range validateIDTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		result, _ := validateID(tc.id, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateFormat(t *testing.T) {
	for _, tc := range validateFormatTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		result, _ := validateFormat(tc.format, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateClass(t *testing.T) {
	for _, tc := range validateClassTC {
		r := NewHtsgetRequest()
		result, _ := validateClass(tc.class, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateReferenceName(t *testing.T) {
	for _, tc := range validateReferenceNameTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		r.AddScalarParam("id", tc.id)
		r.AddScalarParam("class", tc.class)
		result, _ := validateReferenceName(tc.referenceName, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateStart(t *testing.T) {
	for _, tc := range validateStartTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		r.AddScalarParam("referenceName", tc.referenceName)
		result, _ := validateStart(tc.start, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateEnd(t *testing.T) {
	for _, tc := range validateEndTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		r.AddScalarParam("referenceName", tc.referenceName)
		r.AddScalarParam("start", tc.start)
		result, _ := validateEnd(tc.end, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateFields(t *testing.T) {
	for _, tc := range validateFieldsTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		result, _ := validateFields(tc.fields, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateTags(t *testing.T) {
	for _, tc := range validateTagsTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		result, _ := validateTags(tc.tags, r)
		assert.Equal(t, tc.exp, result)
	}
}

func TestValidateNoTags(t *testing.T) {
	for _, tc := range validateNoTagsTC {
		r := NewHtsgetRequest()
		r.AddScalarParam("class", tc.class)
		r.AddListParam("tags", tc.tags)
		result, _ := validateNoTags(tc.notags, r)
		assert.Equal(t, tc.exp, result)
	}
}
