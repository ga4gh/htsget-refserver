// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module samtools view defines the job submission for the 'samtools view'
// command, which streams SAM or BAM to stdout
package htscli

import (
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

// SamtoolsViewCommand represents a single 'samtools view' command and associated
// arguments
type SamtoolsViewCommand struct {
	command *Command
}

// SamtoolsView instantiates a new SamtoolsView Command
func SamtoolsView() *SamtoolsViewCommand {
	// consistent base command and first arg
	samtoolsViewCommand := new(SamtoolsViewCommand)
	samtoolsViewCommand.command = NewCommand()
	samtoolsViewCommand.command.SetBaseCommand("samtools")
	samtoolsViewCommand.command.AddArg("view")
	return samtoolsViewCommand
}

// AddFilePath adds a url or local path to the input file as a cli arg
func (samtoolsViewCommand *SamtoolsViewCommand) AddFilePath(filepath string) *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg(filepath)
	return samtoolsViewCommand
}

// HeaderIncluded adds an option to the cli, which will lead to the header being
// added to the output
func (samtoolsViewCommand *SamtoolsViewCommand) HeaderIncluded() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-h")
	return samtoolsViewCommand
}

// HeaderOnly adds an option to the cli, which will lead to only the header being
// printed to stdout
func (samtoolsViewCommand *SamtoolsViewCommand) HeaderOnly() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-H")
	return samtoolsViewCommand
}

// OutputBAM adds an option to the cli, which will lead to the output being printed
// in BAM format instead of SAM
func (samtoolsViewCommand *SamtoolsViewCommand) OutputBAM() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-b")
	return samtoolsViewCommand
}

// AddRegion adds a specific region request to the command line
func (samtoolsViewCommand *SamtoolsViewCommand) AddRegion(region *htsrequest.Region) *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg(region.ExportSamtools())
	return samtoolsViewCommand
}

// StreamFromStdin adds a cli option, indicating that the input will come from
// stdin and not an input file
func (samtoolsViewCommand *SamtoolsViewCommand) StreamFromStdin() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-")
	return samtoolsViewCommand
}

// GetCommand exports the SamtoolsViewCommand as a generic Command
func (samtoolsViewCommand *SamtoolsViewCommand) GetCommand() *Command {
	return samtoolsViewCommand.command
}
