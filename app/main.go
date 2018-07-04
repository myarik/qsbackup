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
				if _, err = b.DB.Create(dirPath, dirHash, archive); err != nil {
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
						err = b.Storage.Delete(lastArchive, b.Logger)
						if err != nil {
							b.Logger.Error(
								fmt.Sprintf("can't delete the archive %s: %s\n", lastArchive.Name, err))
							return
						}
						b.Logger.Info(fmt.Sprintf("%s backup deleted", lastArchive.Name))
					}
					backups = backups[len(backups)-numberBackups:]
					err = b.DB.Update(dirPath, backups)
					if err != nil {
						b.Logger.Error(
							fmt.Sprintf("can't update the db record %s: %s\n", dirPath, err))
						return
					}
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

// ShowLastBackups returns a dir's last backup
func (b *Backup) ShowLastBackups() {
	for _, item := range b.BackupDirs {
		lastBackup, err := b.DB.Last(item)
		if err != nil {
			b.Logger.Error(fmt.Sprintf("can't get a last backup, for a dir %s: %s", item, err))
			return
		}
		if lastBackup != nil {
			fmt.Printf("The %s was backuped on %s\n",
				item, lastBackup.Timestamp.Format(time.RFC1123))
		} else {
			fmt.Printf("The %s hasn't backuped yet\n", item)
		}
	}
}

// AllBackups returns a dir backups
func (b *Backup) AllBackups() {
	for _, item := range b.BackupDirs {
		listBackups, err := b.DB.List(item)
		if err != nil {
			b.Logger.Error(fmt.Sprintf("can't get a last backup, for a dir %s: %s", item, err))
			return
		}
		if len(listBackups) == 0 {
			fmt.Printf("The %s hasn't backuped yet\n", item)
		} else {
			fmt.Printf("The dir %s:\n", item)
			for _, backupInfo := range listBackups {
				fmt.Printf("The %s was backuped on %s: %s\n",
					item, backupInfo.Timestamp.Format(time.RFC1123), backupInfo.BackupPath)
			}
		}
	}
}
