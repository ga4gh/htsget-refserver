package htscli

import (
	"github.com/ga4gh/htsget-refserver/internal/htsrequest"
)

type SamtoolsViewCommand struct {
	command *Command
}

func SamtoolsView() *SamtoolsViewCommand {
	samtoolsViewCommand := new(SamtoolsViewCommand)
	samtoolsViewCommand.command = NewCommand()
	samtoolsViewCommand.command.SetBaseCommand("samtools")
	samtoolsViewCommand.command.AddArg("view")
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) AddFilePath(filepath string) *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg(filepath)
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) HeaderIncluded() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-h")
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) HeaderOnly() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-H")
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) OutputBAM() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-b")
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) AddRegion(region *htsrequest.Region) *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg(region.ExportSamtools())
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) StreamFromStdin() *SamtoolsViewCommand {
	samtoolsViewCommand.command.AddArg("-")
	return samtoolsViewCommand
}

func (samtoolsViewCommand *SamtoolsViewCommand) GetCommand() *Command {
	return samtoolsViewCommand.command
}
