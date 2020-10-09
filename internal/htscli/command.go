package htscli

import (
	"os/exec"
)

type Command struct {
	baseCommand string
	args        []string
	cmd         *exec.Cmd
}

func NewCommand() *Command {
	command := new(Command)
	command.args = []string{}
	return command
}

func (command *Command) SetBaseCommand(baseCommand string) {
	command.baseCommand = baseCommand
}

func (command *Command) SetArgs(args []string) {
	command.args = args
}

func (command *Command) AddArg(arg string) {
	command.args = append(command.args, arg)
}

func (command *Command) SetupCmd() {
	cmd := exec.Command(command.baseCommand, command.args...)
	command.cmd = cmd
}

func (command *Command) ExecuteCmd() {
	command.cmd.Start()
}
