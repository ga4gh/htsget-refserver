// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module cliargs contains operations for setting properties from the
// command line
package htsconfig

import (
	"flag"
	"sync"
)

// cliArgs contains all properties that can be specified on the command line
type cliArgs struct {
	configFile string
}

// cliargs (*cliArgs): singleton of settings loaded from command line
var cliargs *cliArgs

// cliargsLoaded (sync.Once): indicates whether the singleton has been loaded or not
var cliargsLoaded sync.Once

// parseCliArgs parse all cli options/flags and returns it as a new cliArgs
// instance
func parseCliArgs() *cliArgs {
	configFilePtr := flag.String("config", "", "path to json config file")
	flag.Parse()
	newCliargs := new(cliArgs)
	newCliargs.configFile = *configFilePtr
	return newCliargs
}

// loadCliArgs instantiates the cliargs config singleton, loading allowed options
// into the object
func loadCliArgs() {
	cliargs = parseCliArgs()
}

// getCliArgs gets the loaded cliargs singleton, loading it first if it hasn't
// already been loaded
func getCliArgs() *cliArgs {
	cliargsLoaded.Do(func() {
		loadCliArgs()
	})
	return cliargs
}
