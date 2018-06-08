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
	"strings"
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
	assert.True(t, strings.HasPrefix(res.ID, "3359227887"))
}

func TestBoltDB_madeID(t *testing.T) {
	defer os.Remove(testDB)
	db := prepData()
	res := db.makeID("/home/test")
	assert.True(t, strings.HasPrefix(res, "43430298"))
}

func TestBoltDB_List(t *testing.T) {
	defer os.Remove(testDB)
	db := prepData()
	db.Create("/home/test/", "asdd", "/backup/test.zip")
	db.Create("/home/test2/", "asdd", "/backup/test.zip")
	db.Create("/home/test/", "asdd", "/backup/test.zip")
	res, err := db.List("/home/test/")
	require.NoError(t, err)
	assert.Equal(t, len(res), 2)
	res, err = db.List("/home/test2/")
	require.NoError(t, err)
	assert.Equal(t, len(res), 1)
}

func TestBoltDB_Last(t *testing.T) {
	defer os.Remove(testDB)
	db := prepData()
	res, err := db.Last("/home/test/")
	assert.Nil(t, res)
	assert.Nil(t, err)
	db.Create("/home/test/", "asdd", "/backup/test.zip")
	db.Create("/home/test/", "asdd", "/backup/test2.zip")
	res, _ = db.Last("/home/test/")
	assert.Equal(t, res.BackupPath, "/backup/test2.zip")
}

func TestBoltDB_Pop(t *testing.T) {
	defer os.Remove(testDB)
	db := prepData()
	db.Create("/home/test/", "asdd", "/backup/test.zip")
	db.Create("/home/test/", "asdd", "/backup/test2.zip")
	db.Pop("/home/test/")
	res, _ := db.List("/home/test/")
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].BackupPath, "/backup/test2.zip")
}

func prepData() *BoltDB {
	dbTest, _ := NewBoltDB(testDB, bolt.Options{Timeout: 30 * time.Second}, log.New(ioutil.Discard, 1))
	return dbTest
}
