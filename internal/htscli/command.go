// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module command contains definitions for a single command-line job
package htscli

import (
	"os/exec"
)

// Command job/command to be submitted on the command-line
type Command struct {
	baseCommand string
	args        []string
	cmd         *exec.Cmd
}

// NewCommand instantiates a new Command
func NewCommand() *Command {
	command := new(Command)
	command.args = []string{}
	return command
}

// SetBaseCommand sets the command's base command, ie. the first command
// as it appears on the command line
func (command *Command) SetBaseCommand(baseCommand string) {
	command.baseCommand = baseCommand
}

// SetArgs sets the command's arguments, ie. the space-delimited strings
// appearing after the base command to modify program behaviour
func (command *Command) SetArgs(args []string) {
	command.args = args
}

// AddArg adds a single argument to the arg array
func (command *Command) AddArg(arg string) {
	command.args = append(command.args, arg)
}

// GetArgs gets the argument array of a command
func (command *Command) GetArgs() []string {
	return command.args
}

// GetLastArg gets the final argument in the array
func (command *Command) GetLastArg() string {
	return command.args[len(command.args)-1]
}

// SetupCmd wraps the command's base command and arguments as an exec.Cmd
// object, setting it to the command's cmd property
func (command *Command) SetupCmd() {
	cmd := exec.Command(command.baseCommand, command.args...)
	command.cmd = cmd
}

// ExecuteCmd starts the command that has been set up
func (command *Command) ExecuteCmd() {
	command.cmd.Start()
}
