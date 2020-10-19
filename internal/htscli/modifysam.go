// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module modifysam defines the job submission for the 'modify-sam' subcommand
// of the 'htsget-refserver-utils' cli app, which allows custom
// inclusion/exclusion of SAM fields and tags
package htscli

import (
	"strings"
)

// ModifySamCommand represents the 'modify-sam' subcommand of the
// 'htsget-refserver-utils' package
type ModifySamCommand struct {
	fields []string
	tags   []string
	notags []string
}

// ModifySam instantiates a new ModifySam command
func ModifySam() *ModifySamCommand {
	modifySamCommand := new(ModifySamCommand)
	modifySamCommand.fields = []string{}
	modifySamCommand.tags = []string{}
	modifySamCommand.notags = []string{}
	return modifySamCommand
}

// SetFields sets the ModifySam's '-fields' property
func (modifySamCommand *ModifySamCommand) SetFields(fields []string) {
	modifySamCommand.fields = fields
}

// SetTags sets the ModifySam's '-tags' property
func (modifySamCommand *ModifySamCommand) SetTags(tags []string) {
	modifySamCommand.tags = tags
}

// SetNoTags sets the ModifySam's '-notags' property
func (modifySamCommand *ModifySamCommand) SetNoTags(notags []string) {
	modifySamCommand.notags = notags
}

// GetCommand exports the contents of the ModifySam object as a Command
func (modifySamCommand *ModifySamCommand) GetCommand() *Command {
	// consistent base command and first arg
	command := new(Command)
	command.SetBaseCommand("htsget-refserver-utils")
	command.AddArg("modify-sam")

	// add fields cli flag and value to args list
	if len(modifySamCommand.fields) > 0 {
		command.AddArg("-fields")
		command.AddArg(strings.Join(modifySamCommand.fields, ","))
	}

	// add tags cli flag and value to args list
	if len(modifySamCommand.tags) > 0 {
		command.AddArg("-tags")
		command.AddArg(strings.Join(modifySamCommand.tags, ","))
	}

	// add notags cli flag and value to args list
	if len(modifySamCommand.notags) > 0 {
		command.AddArg("-notags")
		command.AddArg(strings.Join(modifySamCommand.notags, ","))
	}

	return command
}
