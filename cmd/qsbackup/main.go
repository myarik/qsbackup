package main

import (
	pkgLogger "github.com/myarik/qsbackup/pkg/logger"
	"os"
	"github.com/myarik/qsbackup/pkg/config"
	"io/ioutil"
	"log"
)

func main() {
	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Can't open configuration file: %s", filename)
	}
	conf := config.Load(source)
	logger := pkgLogger.New(os.Stdout, 0)
	logger.Debug(conf.Name)
}
