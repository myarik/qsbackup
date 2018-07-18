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

// backupDir creates a dir backup
func (b *Backup) backupDir(dirPath string) {

	lastBackup, err := b.DB.Last(dirPath)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("can't get a last backup, for a dir %s: %s\n", dirPath, err))
		return
	}
	// Create a dir hash
	dirHash, err := store.DirHash(dirPath)
	if err != nil {
		b.Logger.Error(
			fmt.Sprintf("can't build a dir hash, for a dir %s: %s\n", dirPath, err))
		return
	}

	if lastBackup == nil || lastBackup.Hash != dirHash {
		// Create a backup
		archive, err := b.Storage.Save(dirPath, b.Logger)
		if err != nil {
			return
		}
		// Create a db record
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

		// Remove old files if we have more than 4 copies
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
}

// Run is runner
func (b *Backup) Run(jobs int32, cancelChan <-chan struct{}) (<-chan struct{}, error) {
	taskStopChan := make(chan struct{}, 1)
	go func() {
		defer func() {
			taskStopChan <- struct{}{}
		}()
		limit := make(chan bool, jobs)
		var wg sync.WaitGroup

		for _, backupDir := range b.BackupDirs {
			wg.Add(1)
			go func(dirPath string) {
				limit <- true
				defer func() {
					<-limit
					wg.Done()
				}()
				select {
				case <-cancelChan:
					b.Logger.Warning(fmt.Sprintf("%s was canceled\n", dirPath))
				default:
					b.backupDir(dirPath)
				}
			}(backupDir)
		}
		wg.Wait()
	}()
	return taskStopChan, nil
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
