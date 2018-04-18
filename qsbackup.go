package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"sync"
)

type Backup struct {
	Logger *logger.Log
	BackupDirs []string
	Storage Storage
}

// New creates a new Backup.
func New(Config *BackupConfig, Logger *logger.Log) *Backup {
	var backupDirs []string
	for _, backupDir := range Config.Dirs {
		backupDirs = append(backupDirs, backupDir.Path)
	}
	// TODO Put the dest path to the config
	return &Backup{
		Logger: Logger,
		BackupDirs: backupDirs,
		Storage: &LocalStorage{
			Archiver: ZIP,
			DestPath: "/Users/yaroslavmuravskiy/Documents/test_backup/tmp/",
		},
	}
}

func (b *Backup) Run() error {
	// TODO Put the limit value to the config
	limit := make(chan bool, 3)
	var wg sync.WaitGroup

	for _, backupDir := range b.BackupDirs {
		wg.Add(1)
		go func(path string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			limit <- true
			b.Storage.Save(path, b.Logger)
			<-limit
		}(backupDir)
	}
	wg.Wait()

	return nil
}
