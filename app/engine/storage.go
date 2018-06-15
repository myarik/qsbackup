package engine

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
	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage represents type capable to save an archive
type Storage interface {
	Save(src string, logger *logger.Log) (*Archive, error)
	Delete(archive *Archive, logger *logger.Log) error
}

// LocalStorage represents the local storage
type LocalStorage struct {
	Archiver Archiver
	DestPath string
}

// Save backups a dir and saves it in the dest path
func (storage *LocalStorage) Save(src string, logger *logger.Log) (*Archive, error) {
	if err := os.MkdirAll(filepath.Dir(storage.DestPath), 0644); err != nil {
		return nil, err
	}
	archive, err := storage.Archiver.Archive(src, storage.DestPath, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't archive the directory %s, %s", src, err))
		return nil, err
	}
	return archive, nil
}

// Delete a file
func (storage *LocalStorage) Delete(archive *Archive, logger *logger.Log) error {
	err := os.Remove(archive.ArchivePath)
	if err != nil {
		logger.Error(fmt.Sprintf("can't delete a file, %s, error: %s", archive.ArchivePath, err))
	}
	return err
}

// AwsStorage represents the AWS storage
type AwsStorage struct {
	Archiver     Archiver
	AccessKeyID  string
	AccessSecret string
	Region       string
	AWSS3Bucket  string
}

// Save backups a dir and saves it in the S3
func (storage *AwsStorage) Save(src string, logger *logger.Log) (*Archive, error) {
	tempDir, err := ioutil.TempDir("", "qsbackup")
	if err != nil {
		logger.Error(fmt.Sprintf("Can't create the tmp directory, %s", err))
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Created the tmp directory, %s", tempDir))
	defer os.RemoveAll(tempDir)

	archive, err := storage.Archiver.Archive(src, tempDir, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't archive the directory %s, %s", src, err))
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Archived the %s to %s", src, archive))

	uploader := s3manager.NewUploader(storage.getSession())
	f, err := os.Open(archive.ArchivePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't read a file %s, %v", archive, err))
		return nil, fmt.Errorf("can't read a file %s, %v", archive, err)
	}
	defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("private"),
		Bucket: aws.String(storage.AWSS3Bucket),
		Key:    aws.String(filepath.Base(archive.ArchivePath)),
		Body:   f,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to upload file, %v", err))
		return nil, err
	}
	archive.ArchivePath = result.Location
	return archive, nil
}

// Delete an object from the S3
func (storage *AwsStorage) Delete(archive *Archive, logger *logger.Log) error {
	svc := s3.New(storage.getSession())

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(storage.AWSS3Bucket),
		Key:    aws.String(archive.ArchiveName),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		logger.Error(fmt.Sprintf("can't delete a file, %s, error: %s", archive.ArchiveName, err))
	}
	return err
}

func (storage *AwsStorage) getSession() *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(storage.Region),
		Credentials: credentials.NewStaticCredentials(
			storage.AccessKeyID,
			storage.AccessSecret,
			""),
	}))
	return sess
}
