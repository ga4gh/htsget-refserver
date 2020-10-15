// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request, and performing validation and
// transformation
//
// Module transformations_test tests module validation
package htsrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// paramTransformer ParamTransformer singleton used for testing
var paramTransformer = NewParamTransformer()

// noTransformTC test cases for NoTransform
var noTransformTC = []struct {
	s, expMsg string
}{
	{"chr1", ""}, {"chr25", ""}, {"chrMT", ""},
}

// transformStringUppercaseTC test cases for TransformStringUppercase
var transformStringUppercaseTC = []struct {
	input, expOutput, expMsg string
}{
	{"bam", "BAM", ""},
	{"cram", "CRAM", ""},
	{"VCF", "VCF", ""},
}

// transformStringLowercaseTC test cases for TransformStringLowercase
var transformStringLowercaseTC = []struct {
	input, expOutput, expMsg string
}{
	{"HEADER", "header", ""},
	{"BAM", "bam", ""},
	{"cram", "cram", ""},
}

// transformStringToIntTC test cases for TransformStringToInt
var transformStringToIntTC = []struct {
	input     string
	expOutput int
	expMsg    string
}{
	{"100000", 100000, ""},
	{"9999999", 9999999, ""},
	{"NaN", 0, "Could not parse value: 'NaN', integer expected"},
}

// transformSplitTC test cases for TransformSplit
var transformSplitTC = []struct {
	input     string
	expOutput []string
	expMsg    string
}{
	{"MD,NI,NH", []string{"MD", "NI", "NH"}, ""},
	{"HZ,AU", []string{"HZ", "AU"}, ""},
	{"LM,LI,OO,AA", []string{"LM", "LI", "OO", "AA"}, ""},
}

// transformSplitAndUppercaseTC test cases for TransformSplitAndUppercase
var transformSplitAndUppercaseTC = []struct {
	input     string
	expOutput []string
	expMsg    string
}{
	{"seq,qual", []string{"SEQ", "QUAL"}, ""},
	{"qname,flag,rname", []string{"QNAME", "FLAG", "RNAME"}, ""},
	{"rnext,pnext,pos,cigar", []string{"RNEXT", "PNEXT", "POS", "CIGAR"}, ""},
}

// TestNoTransform tests NoTransform function
func TestNoTransform(t *testing.T) {
	for _, tc := range noTransformTC {
		result, msg := paramTransformer.NoTransform(tc.s)
		assert.Equal(t, tc.s, result)
		assert.Equal(t, tc.expMsg, msg)
	}
}

// TestTransformStringUppercase tests TransformStringUppercase function
func TestTransformStringUppercase(t *testing.T) {
	for _, tc := range transformStringUppercaseTC {
		output, msg := paramTransformer.TransformStringUppercase(tc.input)
		assert.Equal(t, tc.expOutput, output)
		assert.Equal(t, tc.expMsg, msg)
	}
}

// TestTransformStringLowercase tests TransformStringLowercase function
func TestTransformStringLowercase(t *testing.T) {
	for _, tc := range transformStringLowercaseTC {
		output, msg := paramTransformer.TransformStringLowercase(tc.input)
		assert.Equal(t, tc.expOutput, output)
		assert.Equal(t, tc.expMsg, msg)
	}
}

// TestTransformStringToInt tests TransformStringToInt function
func TestTransformStringToInt(t *testing.T) {
	for _, tc := range transformStringToIntTC {
		output, msg := paramTransformer.TransformStringToInt(tc.input)
		assert.Equal(t, tc.expOutput, output)
		assert.Equal(t, tc.expMsg, msg)
	}
}

// TestTransformSplit tests TransformSplit function
func TestTransformSplit(t *testing.T) {
	for _, tc := range transformSplitTC {
		output, msg := paramTransformer.TransformSplit(tc.input)
		assert.Equal(t, tc.expOutput, output)
		assert.Equal(t, tc.expMsg, msg)
	}
}

// TestTransformSplitAndUppercase tests TransformSplitAndUppercase function
func TestTransformSplitAndUppercase(t *testing.T) {
	for _, tc := range transformSplitAndUppercaseTC {
		output, msg := paramTransformer.TransformSplitAndUppercase(tc.input)
		assert.Equal(t, tc.expOutput, output)
		assert.Equal(t, tc.expMsg, msg)
	}
}
