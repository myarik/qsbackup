package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"path/filepath"
	"os"
	"fmt"
	"errors"
	"archive/zip"
	"io"
	"time"
	"path"
)

type Backup struct {
	Logger *logger.Log
	Config *BackupConfig
}

var BackupFileError = errors.New("can't create a backup file")

func getDestPath(destPath, srcPath, fileExt string) string {
	now := time.Now()
	return path.Join(
		destPath,
		fmt.Sprintf("%s_%02d_%02d_%d.%s",
			filepath.Base(srcPath),  // File name
			now.Day(),
			now.Month(),
			now.Year(),
			fileExt,
		),
	)
}

func createArchive(src, dest string, logger *logger.Log) error {
	// Create a file and setup a zip writer
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	archiveWriter := zip.NewWriter(out)
	defer archiveWriter.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error(fmt.Sprintf("prevent panic by handling failure accessing a path %q: %v\n", src, err))
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil // skip
		}
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		defer input.Close()
		f, err := archiveWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, input)
		if err != nil {
			return err
		}
		logger.Debug(fmt.Sprintf("Zip file: %q", path))
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", src, err)
	}
	return nil
}

func (b *Backup) Run() error {
	// TODO Put the dest path to the dir
	dest := "/Users/yaroslavmuravskiy/Documents/test_backup/tmp/"
	if err := os.MkdirAll(filepath.Dir(dest), 0644); err != nil {
		return err
	}
	for _, backupDir := range b.Config.Dirs {
		outFile := getDestPath(dest, backupDir.Path, "zip")
		createArchive(backupDir.Path, outFile, b.Logger)
	}
	return nil
}
