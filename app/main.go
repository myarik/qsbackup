package app

import (
	"sync"
	"github.com/myarik/qsbackup/pkg/logger"
	"github.com/myarik/qsbackup/app/engine"
	"github.com/myarik/qsbackup/app/store"
	"fmt"
	"time"
)

const (
	numberBackups = 4
)

// Backup archives dirs and save them in a storage
type Backup struct {
	Logger     *logger.Log
	BackupDirs []string
	Storage    engine.Storage
	DB         *store.BoltDB
}

// Run is runner
func (b *Backup) Run(jobs int32) error {
	limit := make(chan bool, jobs)
	var wg sync.WaitGroup

	for _, backupDir := range b.BackupDirs {
		wg.Add(1)
		go func(dirPath string) {
			limit <- true
			//// Decrement the counter when the goroutine completes.
			defer func() {
				<-limit
				wg.Done()
			}()
			lastBackup, err := b.DB.Last(dirPath)
			if err != nil {
				b.Logger.Error(
					fmt.Sprintf("can't get a last backup, for a dir %s: %s\n", dirPath, err))
				return
			}
			dirHash, err := store.DirHash(dirPath)
			if err != nil {
				b.Logger.Error(
					fmt.Sprintf("can't build a dir hash, for a dir %s: %s\n", dirPath, err))
				return
			}
			// TODO bug if backup runs in the same day
			if lastBackup == nil || lastBackup.Hash != dirHash {
				archive, err := b.Storage.Save(dirPath, b.Logger)
				if err != nil {
					return
				}
				if _, err := b.DB.Create(dirPath, dirHash, archive); err != nil {
					b.Logger.Error(fmt.Sprintf("can't create a db record: %s\n", err))
					return
				}
				b.Logger.Info(fmt.Sprintf("%s backup created", dirPath))
				backups, err := b.DB.List(dirPath)
				if err != nil {
					b.Logger.Error(fmt.Sprintf("can't get a db record: %s\n", err))
					return
				}
				if len(backups) > numberBackups {
					for _, item := range backups[:len(backups)-numberBackups] {
						lastArchive := &engine.Archive{
							Name: item.ArchiveName,
							Path: item.BackupPath,
						}
						b.Storage.Delete(lastArchive, b.Logger)
						b.Logger.Info(fmt.Sprintf("%s backup deleted", lastArchive.Name))
					}
					backups = backups[len(backups)-numberBackups:]
					b.DB.Update(dirPath, backups)
				}
			} else {
				b.Logger.Info(fmt.Sprintf("%s has not changed, after %s",
					dirPath, lastBackup.Timestamp.Format(time.RFC1123)))
			}
		}(backupDir)
	}
	wg.Wait()

	return nil
}
