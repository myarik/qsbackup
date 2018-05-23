package store

import (
	"testing"
	"github.com/coreos/bbolt"
	"time"
	"io/ioutil"
	log "github.com/myarik/qsbackup/pkg/logger"
	"github.com/stretchr/testify/require"
	"os"
	"github.com/stretchr/testify/assert"
)

var testDB = "../../test/test-backup.db"


func TestNewBoltDB(t *testing.T) {
	defer os.Remove(testDB)
	_, err := NewBoltDB(testDB, bolt.Options{Timeout: 30 * time.Second}, log.New(ioutil.Discard, 1))
	require.NoError(t, err)
	require.FileExists(t, testDB)
}

func TestBoltDB_Create(t *testing.T) {
	defer os.Remove(testDB)
	db := prepData()
	res, err := db.Create("/home/test/", "asdd", "/backup/test.zip")
	require.NoError(t, err)
	assert.Equal(t, res.ID, "3359227887")
}

func prepData() *BoltDB{
	dbTest, _ := NewBoltDB(testDB, bolt.Options{Timeout: 30 * time.Second}, log.New(ioutil.Discard, 1))
	return dbTest
}
