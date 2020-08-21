// Package htsconfig allows the program to be configured with modifiable
// properties, affecting runtime properties. also contains program constants
//
// Module cliargs.go contains operations for setting properties from the
// command line
package htsconfig

import (
	"flag"
	"sync"
)

// cliArgs contains all properties that can be specified on the command line
//
// Attributes
//	configFile (string): path to JSON config file
type cliArgs struct {
	configFile string
}

// cliargs (*cliArgs): singleton of settings loaded from command line
var cliargs *cliArgs

// cliargsLoaded (sync.Once): indicates whether the singleton has been loaded or not
var cliargsLoaded sync.Once

// loadCliArgs instantiates the cliargs config singleton, loading allowed options
// into the object
func loadCliArgs() {
	// cli opts
	configFilePtr := flag.String("config", "", "path to json config file")
	// parse command line and construct
	flag.Parse()
	newCliargs := new(cliArgs)
	newCliargs.configFile = *configFilePtr
	cliargs = newCliargs
}

// getCliArgs gets the loaded cliargs singleton, loading it first if it hasn't
// already been loaded
//
// Returns
//	(*cliArgs): parse command line arguments
func getCliArgs() *cliArgs {
	cliargsLoaded.Do(func() {
		loadCliArgs()
	})
	return cliargs
}
