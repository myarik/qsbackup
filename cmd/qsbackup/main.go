package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
)

func main()  {
	logger := log.New(os.Stdout, 0)
	logger.Debug("Test")
	logger.Info("Test")
	logger.Warning("Test")
	logger.Error("Test")
	logger.Critical("Aaaa")
	logger.Info("Test")
}