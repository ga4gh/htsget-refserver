// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module commandchain defines a series of commands, running simultaneously
// as a piped chain, in which the stdout of one command feeds into the stdin
// of the next command
package htscli

import (
	"io"
)

// CommandChain series of commands in which the stdout of one command is piped
// into the stdin of the following command
type CommandChain struct {
	commands []*Command
}

// NewCommandChain instantiates a new CommandChain
func NewCommandChain() *CommandChain {
	commandChain := new(CommandChain)
	commandChain.commands = []*Command{}
	return commandChain
}

// SetCommands sets the array/chain of commands
func (commandChain *CommandChain) SetCommands(commands []*Command) {
	commandChain.commands = commands
}

// AddCommand adds a single command to the array chain of commands
func (commandChain *CommandChain) AddCommand(command *Command) {
	commandChain.commands = append(commandChain.commands, command)
}

// SetupCommandChain stages all commands in the array chain as an exec.Cmd
func (commandChain *CommandChain) SetupCommandChain() {
	for _, command := range commandChain.commands {
		command.SetupCmd()
	}
}

// ExecuteCommandChain starts all commands in the chain. Each command in the
// chain has its stdout and stdin configured according to the jobs that appear
// before and after it in the chain. The stdout pipe of the final command
// is returned
func (commandChain *CommandChain) ExecuteCommandChain() io.ReadCloser {

	// start at command 0, end at second to last command
	for i := 0; i < len(commandChain.commands)-1; i++ {

		// get the current command, and the command that appears after it
		// set the next command's stdin to the stdout pipe of the current command
		// start the current command
		current := commandChain.commands[i].cmd
		next := commandChain.commands[i+1].cmd
		pipe, _ := current.StdoutPipe()
		next.Stdin = pipe
		current.Start()
	}

	// for the last command, return its stdout pipe
	last := commandChain.GetLastCommand().cmd
	pipe, _ := last.StdoutPipe()
	last.Start()
	return pipe
}

// GetLastCommand returns the final command in the array chain
func (commandChain *CommandChain) GetLastCommand() *Command {
	return commandChain.commands[len(commandChain.commands)-1]
}
