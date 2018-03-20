package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
	"github.com/myarik/qsbackup/pkg/config"
	"io/ioutil"
	"fmt"
)

func main() {
	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Can't open configuration file: %s\n", filename)
		os.Exit(1)
	}
	conf, err := config.Load(source)
	if err != nil {
		fmt.Printf("%s: %s\n", err, filename)
		os.Exit(1)
	}
	logger := log.New(os.Stdout, 0)
	logger.Debug(conf.Logfile)
}
