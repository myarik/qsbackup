package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"os"
	"archive/zip"
	"path/filepath"
	"fmt"
	"io"
	"time"
	"path"
)

type Archiver interface {
	Archive(src, dest string, logger *logger.Log) error
}

func getDestPath(src, destPath, fileExt string) string {
	now := time.Now()
	dirName := filepath.Base(src)
	return path.Join(
		destPath,
		fmt.Sprintf("%s_%02d_%02d_%d.%s",
			dirName,
			now.Day(),
			now.Month(),
			now.Year(),
			fileExt,
		),
	)
}

type zipper struct {}

func (z *zipper) Archive(src, destPath string, logger *logger.Log) error{
	// Formatting the destination path to the file
	dest := z.getDestPath(src, destPath)

	// Create a file
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
		// Open a file
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
		logger.Error(fmt.Sprintf("Can't zip a dir %q: %v\n", src, err))
	}
	return nil
}

func (z *zipper) getDestPath(src, destPath string) string {
	return getDestPath(src, destPath, "zip")
}

var ZIP Archiver = (*zipper)(nil)
