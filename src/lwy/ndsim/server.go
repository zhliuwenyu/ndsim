package ndsim

import (
	"flag"
	"os"
)

var configFile string
var cmdHelp bool

func parseFlag() {
	flagSet := flag.NewFlagSet("ndsim", flag.PanicOnError)
	flagSet.BoolVar(&cmdHelp, "help", false, "list all cmd")
	flagSet.StringVar(&configFile, "c", "./conf/ndsim.ini", "conf file path")
	var cmdV bool
	flagSet.BoolVar(&cmdV, "v", false, "for go test cmd")
	defer func() {
		if err := recover(); err != nil {
			os.Exit(0)
		}
	}()
	flagSet.Parse(os.Args[1:])

	if cmdHelp == true {
		flagSet.Usage()
		os.Exit(0)
	}
}

//Run start ndsim server
func Run() {

	if err := LoadConfigFile(configFile); err != nil {
		os.Exit(1)
	}
	if err := initLog(); err != nil {
		os.Exit(2)
	}
}
