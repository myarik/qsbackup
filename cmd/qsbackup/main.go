package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
	"github.com/myarik/qsbackup/app"
	"io/ioutil"
	"fmt"
	"flag"
)

type Options struct {
	ConfigFile     string
	Debug, Version bool
}

var (
	cmdOptions Options
)

const defaultVersion = "0.0.1"

func init() {
	flag.StringVar(&cmdOptions.ConfigFile, "config", "", "Path to the config file")
	flag.StringVar(&cmdOptions.ConfigFile, "c", "", "Path to the config file")

	flag.BoolVar(&cmdOptions.Version, "v", false, "Show the program version")
	flag.BoolVar(&cmdOptions.Debug, "debug", false, "Debug mode")
}

func main() {
	// Parse the CommandLine
	flag.Parse()
	if cmdOptions.Version == true {
		fmt.Printf("Version %s\n", defaultVersion)
		os.Exit(0)
	}
	if cmdOptions.ConfigFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read and validate a config file
	source, err := ioutil.ReadFile(cmdOptions.ConfigFile)
	if err != nil {
		fmt.Printf("Can't open configuration file: %s\n", cmdOptions.ConfigFile)
		os.Exit(1)
	}
	conf, err := app.ConfigLoad(source)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Setup logger
	logger, err := log.Init(conf.Logfile, cmdOptions.Debug)
	defer logger.Close()
	if err != nil {
		fmt.Printf("Can't setup a logger: %s\n", err)
		os.Exit(1)
	}
	backup, err := app.New(conf, logger)
	if err != nil {
		fmt.Printf("Can't create a runner: %s\n", err)
		os.Exit(1)
	}
	backup.Run()
}
