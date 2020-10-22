// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module bcftoolsview_test tests module bcftoolsview
package htscli

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"

	"github.com/stretchr/testify/assert"
)

// bcftoolsViewSetFilePathTC test cases for SetFilePath
var bcftoolsViewSetFilePathTC = []struct {
	filepath string
}{
	{"/path/to/the/file"},
	{"https://genomics.com/datasets/object0001"},
	{"ftp://datasources.org/genomes/12345"},
}

// bcftoolsViewSetHeaderOnlyTC test cases for SetHeaderOnly
var bcftoolsViewSetHeaderOnlyTC = []struct {
	headerOnly bool
}{
	{true},
	{false},
}

// bcftoolsViewSetRegionTC test cases for SetRegion
var bcftoolsViewSetRegionTC = []struct {
	region *htsrequest.Region
}{
	{
		&htsrequest.Region{
			ReferenceName: "chr1",
			Start:         intPtr(2000000),
			End:           intPtr(3000000),
		},
	},
	{
		&htsrequest.Region{
			ReferenceName: "chr22",
			Start:         intPtr(600000),
			End:           intPtr(999999),
		},
	},
}

// bcftoolsViewGetCommandTC test cases for GetCommand
var bcftoolsViewGetCommandTC = []struct {
	filepath   string
	headerOnly bool
	region     *htsrequest.Region
	expArgs    []string
}{
	{
		"/path/to/the/file",
		true,
		nil,
		[]string{"view", "/path/to/the/file", "--no-version", "-h", "-O", "v"},
	},
	{
		"https://genomics.com/datasets/object0001",
		false,
		&htsrequest.Region{
			ReferenceName: "chr1",
			Start:         intPtr(2000000),
			End:           intPtr(3000000),
		},
		[]string{"view", "https://genomics.com/datasets/object0001", "--no-version",
			"-H", "-O", "v", "-r", "chr1:2000000-3000000"},
	},
}

// TestBcftoolsViewSetFilePath tests SetFilePath function
func TestBcftoolsViewSetFilePath(t *testing.T) {
	for _, tc := range bcftoolsViewSetFilePathTC {
		bcftoolsView := BcftoolsView()
		bcftoolsView.SetFilePath(tc.filepath)
		assert.Equal(t, bcftoolsView.filePath, tc.filepath)
	}
}

// TestBcftoolsViewSetHeaderOnly tests SetHeaderOnly function
func TestBcftoolsViewSetHeaderOnly(t *testing.T) {
	for _, tc := range bcftoolsViewSetHeaderOnlyTC {
		bcftoolsView := BcftoolsView()
		bcftoolsView.SetHeaderOnly(tc.headerOnly)
		assert.Equal(t, bcftoolsView.headerOnly, tc.headerOnly)
	}
}

// TestBcftoolsViewSetRegion tests SetRegion function
func TestBcftoolsViewSetRegion(t *testing.T) {
	for _, tc := range bcftoolsViewSetRegionTC {
		bcftoolsView := BcftoolsView()
		bcftoolsView.SetRegion(tc.region)
		assert.Equal(t, bcftoolsView.region.String(), tc.region.String())
	}
}

// TestBcftoolsViewGetCommand tests GetCommand function
func TestBcftoolsViewGetCommand(t *testing.T) {
	for _, tc := range bcftoolsViewGetCommandTC {
		bcftoolsView := BcftoolsView()
		bcftoolsView.SetFilePath(tc.filepath)
		bcftoolsView.SetHeaderOnly(tc.headerOnly)
		bcftoolsView.SetRegion(tc.region)
		command := bcftoolsView.GetCommand()
		assert.Equal(t, "bcftools", command.baseCommand)
		assert.Equal(t, tc.expArgs, command.GetArgs())
	}
}
