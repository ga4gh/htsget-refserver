// Package htscli deals with the construction and submission of command-line
// jobs
//
// Module commandchain_test tests module commandchain
package htscli

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// commandChainSetCommandTC test cases for SetCommands
var commandChainSetCommandsTC = []struct {
	commands []*Command
}{
	{
		[]*Command{
			&Command{
				baseCommand: "echo",
				args:        []string{"Hello", "World"},
			},
			&Command{
				baseCommand: "ls",
				args:        []string{"-1"},
			},
			&Command{
				baseCommand: "wc",
				args:        []string{"-l"},
			},
		},
	},
}

// commandChainExecuteCommandChainTC test cases for ExecuteCommandChain
var commandChainExecuteCommandChainTC = []struct {
	commands  []*Command
	expStdout string
}{
	{
		[]*Command{
			&Command{
				baseCommand: "echo",
				args:        []string{"Hello", "World"},
			},
		},
		"Hello World\n",
	},
	{
		[]*Command{
			&Command{
				baseCommand: "echo",
				args:        []string{"-e", "\"A\nB\nC\nD\nE\""},
			},
			&Command{
				baseCommand: "cat",
				args:        []string{},
			},
			&Command{
				baseCommand: "wc",
				args:        []string{"-l"},
			},
		},
		"       5\n",
	},
}

// commandChainAddCommandTC test cases for AddCommand
var commandChainAddCommandTC = commandChainSetCommandsTC

// TestCommandChainSetCommands tests SetCommands function
func TestCommandChainSetCommands(t *testing.T) {
	for _, tc := range commandChainSetCommandsTC {
		commandChain := NewCommandChain()
		commandChain.SetCommands(tc.commands)
		for c := 0; c < len(tc.commands); c++ {
			expCommand := tc.commands[c]
			actualCommand := commandChain.commands[c]
			assert.Equal(t, expCommand.baseCommand, actualCommand.baseCommand)

			for i := 0; i < len(expCommand.args); i++ {
				expArg := expCommand.args[i]
				actualArg := actualCommand.args[i]
				assert.Equal(t, expArg, actualArg)
			}
		}
	}
}

// TestCommandChainAddCommand tests AddCommand function
func TestCommandChainAddCommand(t *testing.T) {
	for _, tc := range commandChainAddCommandTC {
		commandChain := NewCommandChain()
		for i := 0; i < len(tc.commands); i++ {
			commandChain.AddCommand(tc.commands[i])
		}

		for c := 0; c < len(tc.commands); c++ {
			expCommand := tc.commands[c]
			actualCommand := commandChain.commands[c]
			assert.Equal(t, expCommand.baseCommand, actualCommand.baseCommand)

			for i := 0; i < len(expCommand.args); i++ {
				expArg := expCommand.args[i]
				actualArg := actualCommand.args[i]
				assert.Equal(t, expArg, actualArg)
			}
		}
	}
}

// TestCommandChainExecuteCommandChain tests ExecuteCommandChain function
func TestCommandChainExecuteCommandChain(t *testing.T) {
	for _, tc := range commandChainExecuteCommandChainTC {
		commandChain := NewCommandChain()
		commandChain.SetCommands(tc.commands)
		commandChain.SetupCommandChain()
		pipe := commandChain.ExecuteCommandChain()
		bytes, err := ioutil.ReadAll(pipe)
		assert.Nil(t, err)

		stdout := string(bytes)
		assert.Equal(t, tc.expStdout, stdout)
	}

}
