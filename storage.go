package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(src string, logger *logger.Log) error
}

type LocalStorage struct {
	Archiver Archiver
	DestPath string
}

func (storage *LocalStorage) Save(src string, logger *logger.Log) error{
	if err := os.MkdirAll(filepath.Dir(storage.DestPath), 0644); err != nil {
		return err
	}
	return storage.Archiver.Archive(src, storage.DestPath, logger)
}
