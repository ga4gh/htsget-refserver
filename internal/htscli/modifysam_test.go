// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module modifysam_test tests module modifysam
package htscli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// modifySamSetFieldsTC test cases for SetFields
var modifySamSetFieldsTC = []struct {
	fields []string
}{
	{[]string{"QNAME", "SEQ", "QUAL"}},
	{[]string{"TLEN", "RNEXT"}},
	{[]string{"RNAME", "FLAG", "RNEXT", "PNEXT"}},
}

// modifySamSetTagsTC test cases for SetTags
var modifySamSetTagsTC = []struct {
	tags []string
}{
	{[]string{"MD", "HI"}},
	{[]string{"NM", "NZ", "AU", "FA"}},
	{[]string{"MD", "MA", "MN"}},
}

// modifySamSetNoTagsTC test cases for SetNoTags
var modifySamSetNoTagsTC = []struct {
	notags []string
}{
	{[]string{"MD", "HI"}},
	{[]string{"NM", "NZ", "AU", "FA"}},
	{[]string{"MD", "MA", "MN"}},
}

// modifySamGetCommandTC test cases for GetCommand
var modifySamGetCommandTC = []struct {
	fields, tags, notags, expArgs []string
}{
	{
		[]string{"QNAME", "FLAG"},
		[]string{"MD", "HI", "HN", "NM"},
		[]string{},
		[]string{"modify-sam", "-fields", "QNAME,FLAG", "-tags", "MD,HI,HN,NM"},
	},
	{
		[]string{},
		[]string{"MD", "HI"},
		[]string{"HU"},
		[]string{"modify-sam", "-tags", "MD,HI", "-notags", "HU"},
	},
	{
		[]string{"QNAME", "RNAME", "CIGAR", "QUAL", "SEQ"},
		[]string{"HN", "MD", "HI"},
		[]string{"AU"},
		[]string{"modify-sam", "-fields", "QNAME,RNAME,CIGAR,QUAL,SEQ", "-tags", "HN,MD,HI", "-notags", "AU"},
	},
}

// TestModifySamSetFields tests SetFields function
func TestModifySamSetFields(t *testing.T) {
	for _, tc := range modifySamSetFieldsTC {
		modifySam := ModifySam()
		modifySam.SetFields(tc.fields)
		for i := 0; i < len(tc.fields); i++ {
			assert.Equal(t, tc.fields[i], modifySam.fields[i])
		}
	}
}

// TestModifySamSetTags tests SetTags function
func TestModifySamSetTags(t *testing.T) {
	for _, tc := range modifySamSetTagsTC {
		modifySam := ModifySam()
		modifySam.SetTags(tc.tags)
		for i := 0; i < len(tc.tags); i++ {
			assert.Equal(t, tc.tags[i], modifySam.tags[i])
		}
	}
}

// TestModifySamSetNoTags tests SetNoTags function
func TestModifySamSetNoTags(t *testing.T) {
	for _, tc := range modifySamSetNoTagsTC {
		modifySam := ModifySam()
		modifySam.SetNoTags(tc.notags)
		for i := 0; i < len(tc.notags); i++ {
			assert.Equal(t, tc.notags[i], modifySam.notags[i])
		}
	}
}

// TestModifySamGetCommand tests GetCommand function
func TestModifySamGetCommand(t *testing.T) {
	for _, tc := range modifySamGetCommandTC {
		modifySam := ModifySam()
		modifySam.SetFields(tc.fields)
		modifySam.SetTags(tc.tags)
		modifySam.SetNoTags(tc.notags)
		command := modifySam.GetCommand()
		assert.Equal(t, "htsget-refserver-utils", command.baseCommand)
		for i := 0; i < len(tc.expArgs); i++ {
			assert.Equal(t, tc.expArgs[i], command.args[i])
		}
	}
}
