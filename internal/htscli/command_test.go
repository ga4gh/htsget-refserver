// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module command_test tests module command
package htscli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// commandSetBaseCommandTC test cases for SetBaseCommand
var commandSetBaseCommandTC = []struct {
	baseCommand string
}{
	{"ls"}, {"pwd"}, {"wc"},
}

// commandSetArgsTC test cases for SetArgs
var commandSetArgsTC = []struct {
	args []string
}{
	{[]string{"-b", "-c", "-l"}},
	{[]string{"--arg", "--config"}},
	{[]string{"-h"}},
}

// commandAddArgTC test cases for AddArg
var commandAddArgTC = []struct {
	args []string
}{
	{[]string{"-b", "-c", "-l"}},
	{[]string{"--arg", "--config"}},
	{[]string{"-h"}},
}

// commandSetupCmdTC test cases for SetupCmd
var commandSetupCmdTC = []struct {
	baseCommand string
	args        []string
}{
	{"echo", []string{"Hello", "World"}},
}

// TestCommandSetBaseCommand tests SetBaseCommand function
func TestCommandSetBaseCommand(t *testing.T) {
	for _, tc := range commandSetBaseCommandTC {
		command := NewCommand()
		command.SetBaseCommand(tc.baseCommand)
		assert.Equal(t, tc.baseCommand, command.baseCommand)
	}
}

// TestCommandSetArgs tests SetArgs function
func TestCommandSetArgs(t *testing.T) {
	for _, tc := range commandSetArgsTC {
		command := NewCommand()
		command.SetArgs(tc.args)
		for i := 0; i < len(tc.args); i++ {
			assert.Equal(t, tc.args[i], command.GetArgs()[i])
		}
		assert.Equal(t, tc.args[len(tc.args)-1], command.GetLastArg())
	}
}

// TestCommandAddArg tests AddArg function
func TestCommandAddArg(t *testing.T) {
	for _, tc := range commandAddArgTC {
		command := NewCommand()
		for i := 0; i < len(tc.args); i++ {
			command.AddArg(tc.args[i])
		}
		for i := 0; i < len(tc.args); i++ {
			assert.Equal(t, tc.args[i], command.GetArgs()[i])
		}
		assert.Equal(t, tc.args[len(tc.args)-1], command.GetLastArg())
	}
}

// TestCommandSetupCmd tests SetupCmd function
func TestCommandSetupCmd(t *testing.T) {
	for _, tc := range commandSetupCmdTC {
		command := NewCommand()
		command.SetBaseCommand(tc.baseCommand)
		command.SetArgs(tc.args)
		command.SetupCmd()

		expArgs := []string{tc.baseCommand}
		expArgs = append(expArgs, tc.args...)
		actualArgs := command.cmd.Args

		for i := 0; i < len(expArgs); i++ {
			assert.Equal(t, expArgs[i], actualArgs[i])
		}
		command.ExecuteCmd()
	}
}
