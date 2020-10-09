package htscli

import (
	"fmt"
	"io"
)

type CommandChain struct {
	commands []*Command
}

func NewCommandChain() *CommandChain {
	commandChain := new(CommandChain)
	commandChain.commands = []*Command{}
	return commandChain
}

func (commandChain *CommandChain) SetCommands(commands []*Command) {
	commandChain.commands = commands
}

func (commandChain *CommandChain) AddCommand(command *Command) {
	commandChain.commands = append(commandChain.commands, command)
}

func (commandChain *CommandChain) SetupCommandChain() {
	for _, command := range commandChain.commands {
		command.SetupCmd()
	}
}

func (commandChain *CommandChain) ExecuteCommandChain() io.ReadCloser {

	for i := 0; i < len(commandChain.commands)-1; i++ {

		fmt.Println(commandChain.commands[i].args)
		fmt.Println("-")

		current := commandChain.commands[i].cmd
		next := commandChain.commands[i+1].cmd
		pipe, _ := current.StdoutPipe()
		next.Stdin = pipe
		current.Start()
	}
	last := commandChain.GetLastCommand().cmd
	pipe, _ := last.StdoutPipe()
	last.Start()
	return pipe
}

func (commandChain *CommandChain) GetLastCommand() *Command {
	return commandChain.commands[len(commandChain.commands)-1]
}
