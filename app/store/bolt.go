package store

import (
	"github.com/coreos/bbolt"
	"time"
	"github.com/myarik/qsbackup/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
)

const dirsBucket = "dirs"

// BoltDB implements store.
type BoltDB struct {
	db *bolt.DB
}

type DirBackup struct {
	ID         string
	SrcPath    string
	BackupPath string
	Hash       string
	Timestamp  time.Time
}

// Create a botldb store
func NewBoltDB(dbPath string, options bolt.Options, logger *logger.Log) (*BoltDB, error) {
	db, err := bolt.Open(dbPath, 0600, &options)
	if err != nil {
		logger.Error(fmt.Sprintf("Can't create a boltdb: %s, %s", dbPath, err))
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(dirsBucket)); e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		logger.Error("failed to create top level bucket")
		return nil, err
	}
	return &BoltDB{db: db}, nil
}

func (b *BoltDB) Create(dirPath, dirHash, location string) (*DirBackup, error) {
	record := &DirBackup{
		ID:         DirID(dirPath),
		SrcPath:    dirPath,
		BackupPath: location,
		Hash:       dirHash,
		Timestamp:  time.Now(),
	}

	jRecord, err := json.Marshal(&record)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal a record")
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dirsBucket))
		err := b.Put([]byte(record.ID), []byte(jRecord))
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create a record in the db")
	}
	return record, nil
}

func (b *BoltDB) Last() {
}

func (b *BoltDB) Delete() {
}

func (b *BoltDB) DeleteAll() {
}

func (b *BoltDB) Count() {
}
