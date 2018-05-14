package qsbackup

import (
	"github.com/myarik/qsbackup/pkg/logger"
	"os"
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"fmt"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Storage represents type capable to save an archive
type Storage interface {
	Save(src string, logger *logger.Log) (string, error)
}

// LocalStorage represents the local storage
type LocalStorage struct {
	Archiver Archiver
	DestPath string
}

// Save backups a dir and saves it in the dest path
func (storage *LocalStorage) Save(src string, logger *logger.Log) (string, error) {
	if err := os.MkdirAll(filepath.Dir(storage.DestPath), 0644); err != nil {
		return "", err
	}
	location, err := storage.Archiver.Archive(src, storage.DestPath, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't archive the directory %s, %s", src, err))
		return "", err
	}
	logger.Info(fmt.Sprintf("The directory %s archived to %s\n", src, location))
	return location, nil
}

// AwsStorage represents the AWS storage
type AwsStorage struct {
	Archiver     Archiver
	AccessKeyID  string
	AccessSecret string
	Region       string
	Bucket       string
}

// Save backups a dir and saves it in the S3
func (storage *AwsStorage) Save(src string, logger *logger.Log) (string, error) {
	tempDir, err := ioutil.TempDir("", "qsbackup")
	if err != nil {
		logger.Error(fmt.Sprintf("Can't create the tmp directory, %s", err))
		return "", err
	}
	logger.Debug(fmt.Sprintf("Created the tmp directory, %s", tempDir))
	defer os.RemoveAll(tempDir)

	archive, err := storage.Archiver.Archive(src, tempDir, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't archive the directory %s, %s", src, err))
		return "", err
	}
	logger.Debug(fmt.Sprintf("Archived the %s to %s", src, archive))

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(storage.Region),
		Credentials: credentials.NewStaticCredentials(
			storage.AccessKeyID,
			storage.AccessSecret,
			""),
	}))
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(archive)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't read a file %s, %v", archive, err))
		return "", fmt.Errorf("can't read a file %s, %v", archive, err)
	}
	defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("private"),
		Bucket: aws.String(storage.Bucket),
		Key:    aws.String(filepath.Base(archive)),
		Body:   f,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to upload file, %v", err))
		return "", err
	}
	logger.Info(fmt.Sprintf("The directory %s archived and uploaded to, %s\n", src, result.Location))
	return result.Location, nil
}
