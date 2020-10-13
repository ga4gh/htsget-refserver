// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module samtoolsview_test tests module samtoolsview
package htscli

import (
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
	"github.com/stretchr/testify/assert"
)

// samtoolsViewAddFilePathTC test cases for AddFilePath
var samtoolsViewAddFilePathTC = []struct {
	filepath string
}{
	{"/path/to/the/file"},
	{"https://genomics.com/datasets/object0001"},
	{"ftp://datasources.org/genomes/12345"},
}

// samtoolsViewAddRegionTC test cases for AddRegion
var samtoolsViewAddRegionTC = []struct {
	region    *htsrequest.Region
	expRegion string
}{
	{
		&htsrequest.Region{
			ReferenceName: "chr1",
			Start:         intPtr(1000000),
			End:           intPtr(2000000),
		},
		"chr1:1000000-2000000",
	},
	{
		&htsrequest.Region{
			ReferenceName: "chr21",
			Start:         intPtr(550000),
			End:           intPtr(990000),
		},
		"chr21:550000-990000",
	},
	{
		&htsrequest.Region{
			ReferenceName: "chrX",
			Start:         intPtr(240000000),
			End:           intPtr(485000000),
		},
		"chrX:240000000-485000000",
	},
}

// samtoolsViewGetCommandTC test cases for GetCommand
var samtoolsViewGetCommandTC = []struct {
	headerIncluded bool
	headerOnly     bool
	outputBam      bool
	useFilePath    bool
	filepath       string
	region         *htsrequest.Region
	expArgs        []string
}{
	{
		false,
		false,
		false,
		true,
		"/path/to/the/file.bam",
		nil,
		[]string{"view", "/path/to/the/file.bam"},
	},
	{
		true,
		false,
		true,
		false,
		"",
		&htsrequest.Region{
			ReferenceName: "chr1",
			Start:         intPtr(1000000),
			End:           intPtr(2000000),
		},
		[]string{"view", "-h", "-b", "-", "chr1:1000000-2000000"},
	},
	{
		false,
		true,
		true,
		false,
		"",
		&htsrequest.Region{
			ReferenceName: "chrX",
			Start:         intPtr(240000000),
			End:           intPtr(485000000),
		},
		[]string{"view", "-H", "-b", "-", "chrX:240000000-485000000"},
	},
}

// intPtr convenience method to get pointer of an int
func intPtr(i int) *int {
	return &i
}

// TestSamtoolsViewAddFilePath tests AddFilePath function
func TestSamtoolsViewAddFilePath(t *testing.T) {
	for _, tc := range samtoolsViewAddFilePathTC {
		samtoolsView := SamtoolsView()
		samtoolsView.AddFilePath(tc.filepath)
		assert.Equal(t, samtoolsView.command.GetLastArg(), tc.filepath)
	}
}

// TestSamtoolsViewHeaderIncluded tests HeaderIncluded function
func TestSamtoolsViewHeaderIncluded(t *testing.T) {
	samtoolsView := SamtoolsView()
	samtoolsView.HeaderIncluded()
	assert.Equal(t, samtoolsView.command.GetLastArg(), "-h")
}

// TestSamtoolsViewHeaderOnly tests HeaderOnly function
func TestSamtoolsViewHeaderOnly(t *testing.T) {
	samtoolsView := SamtoolsView()
	samtoolsView.HeaderOnly()
	assert.Equal(t, samtoolsView.command.GetLastArg(), "-H")
}

// TestSamtoolsViewOutputBAM tests OutputBAM function
func TestSamtoolsViewOutputBAM(t *testing.T) {
	samtoolsView := SamtoolsView()
	samtoolsView.OutputBAM()
	assert.Equal(t, samtoolsView.command.GetLastArg(), "-b")
}

// TestSamtoolsViewAddRegion tests AddRegion function
func TestSamtoolsViewAddRegion(t *testing.T) {
	for _, tc := range samtoolsViewAddRegionTC {
		samtoolsView := SamtoolsView()
		samtoolsView.AddRegion(tc.region)
		assert.Equal(t, samtoolsView.command.GetLastArg(), tc.expRegion)
	}
}

// TestSamtoolsViewStreamFromStdin tests StreamFromStdin function
func TestSamtoolsViewStreamFromStdin(t *testing.T) {
	samtoolsView := SamtoolsView()
	samtoolsView.StreamFromStdin()
	assert.Equal(t, samtoolsView.command.GetLastArg(), "-")
}

// TestSamtoolsViewGetCommand tests GetCommand function
func TestSamtoolsViewGetCommand(t *testing.T) {
	for _, tc := range samtoolsViewGetCommandTC {
		samtoolsView := SamtoolsView()
		if tc.headerIncluded {
			samtoolsView.HeaderIncluded()
		}
		if tc.headerOnly {
			samtoolsView.HeaderOnly()
		}
		if tc.outputBam {
			samtoolsView.OutputBAM()
		}
		if tc.useFilePath {
			samtoolsView.AddFilePath(tc.filepath)
		} else {
			samtoolsView.StreamFromStdin()
		}
		if tc.region != nil {
			samtoolsView.AddRegion(tc.region)
		}

		command := samtoolsView.GetCommand()
		assert.Equal(t, command.baseCommand, "samtools")

		for i := 0; i < len(tc.expArgs); i++ {
			assert.Equal(t, tc.expArgs[i], command.args[i])
		}
	}
}
