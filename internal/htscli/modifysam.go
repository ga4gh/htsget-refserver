package htscli

import (
	"fmt"
	"strings"
)

type ModifySamCommand struct {
	fields []string
	tags   []string
	notags []string
}

func ModifySam() *ModifySamCommand {
	modifySamCommand := new(ModifySamCommand)
	modifySamCommand.fields = []string{}
	modifySamCommand.tags = []string{}
	modifySamCommand.notags = []string{}
	return modifySamCommand
}

func (modifySamCommand *ModifySamCommand) SetFields(fields []string) {
	modifySamCommand.fields = fields
}

func (modifySamCommand *ModifySamCommand) SetTags(tags []string) {
	modifySamCommand.tags = tags
}

func (modifySamCommand *ModifySamCommand) SetNoTags(notags []string) {
	modifySamCommand.notags = notags
}

func (modifySamCommand *ModifySamCommand) GetCommand() *Command {
	command := new(Command)
	command.SetBaseCommand("htsget-refserver-utils")
	command.AddArg("modify-sam")
	if len(modifySamCommand.fields) > 0 {
		command.AddArg("-fields")
		command.AddArg(strings.Join(modifySamCommand.fields, ","))
	}
	if len(modifySamCommand.tags) > 0 {
		fmt.Println("tags has been set!")
		command.AddArg("-tags")
		command.AddArg(strings.Join(modifySamCommand.tags, ","))
	}
	if len(modifySamCommand.notags) > 0 {
		command.AddArg("-notags")
		command.AddArg(strings.Join(modifySamCommand.notags, ","))
	}
	return command
}
