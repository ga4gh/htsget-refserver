package htsgetconfig

import (
	"flag"
	"sync"
)

type cliArgs struct {
	configFile string
}

var cliargs *cliArgs
var cliargsLoaded sync.Once

func loadCliArgs() {
	// cli opts
	configFilePtr := flag.String("config", "", "path to json config file")
	flag.Parse()
	newCliargs := new(cliArgs)
	newCliargs.configFile = *configFilePtr
	cliargs = newCliargs
}

func getCliArgs() *cliArgs {
	cliargsLoaded.Do(func() {
		loadCliArgs()
	})
	return cliargs
}
