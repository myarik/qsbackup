package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
)

func main()  {
	logger := log.New(os.Stdout, 1)
	logger.Debug("Test")
	logger.Info("Test")

}