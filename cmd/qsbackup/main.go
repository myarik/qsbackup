package main

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"os"
	"github.com/myarik/qsbackup/cmd/qsbackup/config"
)

func main()  {

	tc := &config.BackupConfig{
		Name: "Test vasya",
	}

	logger := log.New(os.Stdout, 0)
	logger.Debug(tc.Logfile)
}