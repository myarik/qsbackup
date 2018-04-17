package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"path/filepath"
	"os"
	"errors"
	"sync"
)

type Backup struct {
	Logger *logger.Log
	Config *BackupConfig
}

var BackupFileError = errors.New("can't create a backup file")


func getArchiver(config *BackupConfig) Archiver {
	return ZIP
}

func (b *Backup) Run() error {
	// TODO Put the dest path to the config
	dest := "/Users/yaroslavmuravskiy/Documents/test_backup/tmp/"

	// TODO Put the limit value to the config
	limit := make(chan bool, 3)
	var wg sync.WaitGroup

	archiver := getArchiver(b.Config)

	if err := os.MkdirAll(filepath.Dir(dest), 0644); err != nil {
		return err
	}
	for _, backupDir := range b.Config.Dirs {
		wg.Add(1)
		go func(path string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			limit <- true
			archiver.Archive(path, dest, b.Logger)
			<-limit
		}(backupDir.Path)
	}
	wg.Wait()

	return nil
}
