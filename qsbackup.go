package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"sync"
	"fmt"
)

// Backup archives dirs and save them in a storage
type Backup struct {
	Logger     *logger.Log
	BackupDirs []string
	Storage    Storage
}

// New creates a new Backup.
func New(Config *BackupConfig, Logger *logger.Log) (*Backup, error) {
	var backupDirs []string
	for _, backupDir := range Config.Dirs {
		backupDirs = append(backupDirs, backupDir.Path)
	}
	backup := &Backup{
		Logger:     Logger,
		BackupDirs: backupDirs,
	}
	switch Config.Storage.Type {
	case "local":
		backup.Storage = &LocalStorage{
			Archiver: ZIP,
			DestPath: Config.Storage.DestPath,
		}
	case "aws":
		backup.Storage = &AwsStorage{
			Archiver:     ZIP,
			Region:       Config.Storage.AwsRegion,
			AccessKeyID:  Config.Storage.AwsKey,
			AccessSecret: Config.Storage.AwsSecret,
			Bucket:       Config.Storage.AwsBucket,
		}
	default:
		return nil, fmt.Errorf("storage type does not support")
	}
	return backup, nil
}

// Run is runner
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
