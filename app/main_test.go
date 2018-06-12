package app

import (
	log "github.com/myarik/qsbackup/pkg/logger"
	"io/ioutil"
	"testing"
	"github.com/coreos/bbolt"
	"time"
	"os"
	"github.com/myarik/qsbackup/app/store"
	"github.com/stretchr/testify/assert"
)

type testStorage struct{}

func (storage *testStorage) Save(src string, logger *log.Log) (string, error) {
	return "/tmp/test_path.zip", nil
}

var testDB = "../test/test-backup.db"

func TestBackup_Run(t *testing.T) {
	defer os.Remove(testDB)
	backup := prepStorage()
	values, _ := backup.DB.List("../test/testdata/hash1")
	assert.Equal(t, len(values), 0)
	backup.Run()
	values, _ = backup.DB.List("../test/testdata/hash1")
	assert.Equal(t, len(values), 1)
	backup.Run()
	values, _ = backup.DB.List("../test/testdata/hash1")
	assert.Equal(t, len(values), 1)
	tmpFile, _ := ioutil.TempFile("../test/testdata/hash1/", "test")
	defer os.Remove(tmpFile.Name()) // clean up
	tmpFile.Write([]byte("Test"))
	backup.Run()
	values, _ = backup.DB.List("../test/testdata/hash1")
	assert.Equal(t, len(values), 2)
}

func prepStorage() *Backup {
	return &Backup{
		Logger:     log.New(ioutil.Discard, 1),
		BackupDirs: []string{"../test/testdata/hash1", "../test/testdata/hash2"},
		Storage:    &testStorage{},
		DB:         prepData(),
	}
}

func prepData() *store.BoltDB {
	dbTest, _ := store.NewBoltDB(testDB, bolt.Options{Timeout: 30 * time.Second}, log.New(ioutil.Discard, 1))
	return dbTest
}