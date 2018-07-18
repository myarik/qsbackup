package engine

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"os"
	"archive/zip"
	"path/filepath"
	"fmt"
	"io"
	"time"
	"path"
	"strings"
)

type Archive struct {
	Name string
	Path string
}

// Archiver represents type capable of archiving
type Archiver interface {
	Archive(src, dest string, logger *logger.Log) (*Archive, error)
}

func getDestPath(src, destPath, fileExt string) (string, string) {
	now := time.Now()
	dirName := filepath.Base(src)
	archiveName := fmt.Sprintf("%s_%02d_%02d_%d__%02d_%02d_%02d.%s",
		dirName,
		now.Day(),
		now.Month(),
		now.Year(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		fileExt,
	)
	return archiveName, path.Join(destPath, archiveName)
}

type zipper struct{}

func (z *zipper) Archive(src, destPath string, logger *logger.Log) (*Archive, error) {
	// Formatting the destination path to the file
	archive := &Archive{}
	archive.Name, archive.Path = z.getDestPath(src, destPath)

	// Create a file
	out, err := os.Create(archive.Path)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	archiveWriter := zip.NewWriter(out)
	defer archiveWriter.Close()

	err = filepath.Walk(src, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error(fmt.Sprintf("prevent panic by handling failure accessing a path %q: %v\n", src, err))
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil // skip
		}

		// Skip a symlink
		fi, err := os.Lstat(srcPath)
		if err != nil {
			return nil
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Open a file
		input, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer input.Close()

		archivePath := strings.Replace(srcPath, src, path.Base(src), 1)

		f, err := archiveWriter.Create(archivePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, input)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Can't zip a dir %q: %v\n", src, err))
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Zipped %s to %s", src, archive.Path))
	return archive, nil
}

// return path & filename
func (z *zipper) getDestPath(src, destPath string) (string, string) {
	return getDestPath(src, destPath, "zip")
}

// ZIP is an Archiver that zips files.
var ZIP Archiver = (*zipper)(nil)
