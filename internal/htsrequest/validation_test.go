// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module validation_test tests module validation
package htsrequest

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconstants"
	"github.com/stretchr/testify/assert"
)

// paramValidator ParamValidator singleton used for testing
var paramValidator = NewParamValidator()

// validateIDTC test cases for ValidateID
var validateIDTC = []struct {
	endpoint htsconstants.APIEndpoint
	id       string
	exp      bool
}{
	{htsconstants.APIEndpointReadsTicket, "NoDataSource.00001", false},
	{htsconstants.APIEndpointReadsTicket, "tabulamuris.NoID", false},
	{htsconstants.APIEndpointReadsTicket, "tabulamuris.A1-B000168-3_57_F-1-1_R2", true},
}

// validateFormatTC test cases for validateFormat
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

// validateClassTC test cases for ValidateClass
var validateClassTC = []struct {
	class string
	exp   bool
}{
	{"header", true},
	{"body", false},
	{"otherclass", false},
}

// validateReferenceNameTC test cases for ValidateReferenceName
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
	{
		htsconstants.APIEndpointVariantsTicket,
		"1000genomes.NoID",
		"",
		"chrMT",
		false,
	},
}

// validateStartTC test cases for ValidateStart
var validateStartTC = []struct {
	class, referenceName string
	start                int
	exp                  bool
}{
	{"header", "", 10000, false},
	{"", "*", 10000, false},
	{"", "", 10000, false},
	{"", "chr1", -10000, false},
	{"", "chr1", 10000, true},
	{"", "chr22", 220000, true},
	{"", "chrMT", 3000, true},
}

// validateEndTC test cases for ValidateEnd
var validateEndTC = []struct {
	class, referenceName string
	start, end           int
	exp                  bool
}{
	{"header", "", 10000, 20000, false},
	{"", "*", 10000, 20000, false},
	{"", "", 10000, 20000, false},
	{"", "chr1", 0, -10000, false},
	{"", "chr1", 100000, 50000, false},
	{"", "chr1", 100, 500, true},
}

// validateFieldsTC test cases for ValidateFields
var validateFieldsTC = []struct {
	class  string
	fields []string
	exp    bool
}{
	{"header", []string{}, false},
	{"", []string{"FOO", "FLAG", "QNAME"}, false},
	{"", []string{"FLAG"}, true},
	{"", []string{"TLEN", "SEQ", "QUAL", "FLAG"}, true},
}

// validateTagsTC test cases for ValidateTags
var validateTagsTC = []struct {
	class string
	tags  []string
	exp   bool
}{
	{"header", []string{"NM", "MD"}, false},
	{"", []string{"NM", "MD"}, true},
	{"", []string{"HZ", "MD", "NM", "HI"}, true},
}

// validateNoTagsTC test cases for ValidateNoTags
var validateNoTagsTC = []struct {
	class  string
	tags   []string
	notags []string
	exp    bool
}{
	{"header", []string{"NM", "MD"}, []string{}, false},
	{"", []string{"NM", "MD"}, []string{}, true},
	{"", []string{"NM", "MD"}, []string{"MD"}, false},
}

// validateRegionsTC test cases for ValidateRegions
var validateRegionsTC = []struct {
	endpoint   htsconstants.APIEndpoint
	id         string
	regions    []*Region
	expBool    bool
	expMessage string
}{
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.NoID",
		[]*Region{
			&Region{"chr1", intPointer(200000), intPointer(300000)},
		},
		false,
		"Could not get referenceNames from requested alignment file",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"chr1", intPointer(200000), intPointer(300000)},
		},
		true,
		"",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"chr25", intPointer(200000), intPointer(300000)},
		},
		false,
		"Invalid referenceName in regions list: 'chr25'",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"", intPointer(200000), intPointer(300000)},
		},
		false,
		"Invalid region(s): 'start' cannot be set without 'referenceName'",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"chr1", intPointer(-100), intPointer(300000)},
		},
		false,
		"Invalid region(s): 'start' MUST be greater than or equal to zero",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"", nil, intPointer(300000)},
		},
		false,
		"Invalid region(s): 'end' cannot be set without 'referenceName'",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"chr1", intPointer(200000), intPointer(-100)},
		},
		false,
		"Invalid region(s): 'end' MUST be greater than or equal to zero",
	},
	{
		htsconstants.APIEndpointReadsTicket,
		"tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[]*Region{
			&Region{"chr1", intPointer(300000), intPointer(200000)},
		},
		false,
		"Invalid region(s): 'end' MUST be greater than 'start'",
	},
}

// intPointer convenience method to get pointer of an int
func intPointer(i int) *int {
	return &i
}

// TestNoValidation test NoValidation function
func TestNoValidation(t *testing.T) {
	r := NewHtsgetRequest()
	input := "BAM"
	found, _ := paramValidator.NoValidation(r, input)
	assert.Equal(t, true, found)
}

// TestValidateID tests ValidateID function
func TestValidateID(t *testing.T) {
	for _, tc := range validateIDTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		result, _ := paramValidator.ValidateID(r, tc.id)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateFormat tests ValidateFormat function
func TestValidateFormat(t *testing.T) {
	for _, tc := range validateFormatTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		result, _ := paramValidator.ValidateFormat(r, tc.format)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateClass tests ValidateClass function
func TestValidateClass(t *testing.T) {
	for _, tc := range validateClassTC {
		r := NewHtsgetRequest()
		result, _ := paramValidator.ValidateClass(r, tc.class)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateReferenceName tests ValidateReferenceName function
func TestValidateReferenceName(t *testing.T) {
	for _, tc := range validateReferenceNameTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		r.SetID(tc.id)
		r.SetClass(tc.class)
		result, _ := paramValidator.ValidateReferenceName(r, tc.referenceName)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateStart tests ValidateStart function
func TestValidateStart(t *testing.T) {
	for _, tc := range validateStartTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		r.SetReferenceName(tc.referenceName)
		result, _ := paramValidator.ValidateStart(r, tc.start)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateEnd tests ValidateEnd function
func TestValidateEnd(t *testing.T) {
	for _, tc := range validateEndTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		r.SetReferenceName(tc.referenceName)
		r.SetStart(tc.start)
		result, _ := paramValidator.ValidateEnd(r, tc.end)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateFields tests ValidateFields function
func TestValidateFields(t *testing.T) {
	for _, tc := range validateFieldsTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		result, _ := paramValidator.ValidateFields(r, tc.fields)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateTags tests ValidateTags function
func TestValidateTags(t *testing.T) {
	for _, tc := range validateTagsTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		result, _ := paramValidator.ValidateTags(r, tc.tags)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateNoTags tests ValidateNoTags function
func TestValidateNoTags(t *testing.T) {
	for _, tc := range validateNoTagsTC {
		r := NewHtsgetRequest()
		r.SetClass(tc.class)
		r.SetTags(tc.tags)
		result, _ := paramValidator.ValidateNoTags(r, tc.notags)
		assert.Equal(t, tc.exp, result)
	}
}

// TestValidateRegions tests ValidateRegions function
func TestValidateRegions(t *testing.T) {
	for _, tc := range validateRegionsTC {
		r := NewHtsgetRequest()
		r.SetEndpoint(tc.endpoint)
		r.SetID(tc.id)
		result, message := paramValidator.ValidateRegions(r, tc.regions)
		assert.Equal(t, tc.expBool, result)
		assert.Equal(t, tc.expMessage, message)
	}
}
