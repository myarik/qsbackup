package store

import (
	"github.com/coreos/bbolt"
	"time"
	"github.com/myarik/qsbackup/pkg/logger"
	"fmt"
	"github.com/pkg/errors"
	"encoding/json"
	"github.com/myarik/qsbackup/app/engine"
)

const dirsBucket = "dirs"

// BoltDB implements store.
type BoltDB struct {
	db *bolt.DB
}

type DirBackup struct {
	ID          string
	SrcPath     string
	BackupPath  string
	ArchiveName string
	Hash        string
	Timestamp   time.Time
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

// Create a new record
func (b *BoltDB) Create(dirPath, dirHash string, archive *engine.Archive) (*DirBackup, error) {
	record := DirBackup{
		ID:          b.makeID(dirPath),
		SrcPath:     dirPath,
		BackupPath:  archive.Path,
		ArchiveName: archive.Name,
		Hash:        dirHash,
		Timestamp:   time.Now(),
	}
	var dirBackups []DirBackup

	err := b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(dirsBucket))
		rawValue := c.Get([]byte(record.ID))
		// If record exists the value unmarshal
		if rawValue != nil {
			if e := json.Unmarshal(rawValue, &dirBackups); e != nil {
				return errors.Wrap(e, "failed to unmarshal")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	dirBackups = append(dirBackups, record)
	jRecord, err := json.Marshal(&dirBackups)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal a record")
	}
	if b.save(record.ID, jRecord) != nil {
		return nil, errors.Wrap(err, "can't create a record in the db")
	}
	return &record, nil
}

// List returns list of all dir backup
func (b *BoltDB) List(dirPath string) (list []DirBackup, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(dirsBucket))
		rawValue := c.Get([]byte(b.makeID(dirPath)))
		// If record exists the value unmarshal
		if rawValue != nil {
			if e := json.Unmarshal(rawValue, &list); e != nil {
				return errors.Wrap(e, "failed to unmarshal")
			}
		}
		return nil
	})
	return list, err
}

// Return the last dir backup
func (b *BoltDB) Last(dirPath string) (*DirBackup, error) {
	backups, err := b.List(dirPath)
	if err != nil {
		return nil, err
	}
	if len(backups) > 0 {
		return &backups[len(backups)-1], nil
	}
	return nil, nil
}

// remove a first record
func (b *BoltDB) Pop(dirPath string) error {
	backups, err := b.List(dirPath)
	if err != nil {
		return err
	}
	if len(backups) == 0 {
		return errors.New("empty record")
	}
	backups = backups[1:]
	jRecord, err := json.Marshal(&backups)
	if err != nil {
		return errors.Wrap(err, "can't marshal a record")
	}
	if b.save(b.makeID(dirPath), jRecord) != nil {
		return errors.Wrap(err, "can't create a record in the db")
	}
	return nil
}

// Update a record
func (b *BoltDB) Update(dirPath string, backups []DirBackup) error {
	jRecord, err := json.Marshal(&backups)
	if err != nil {
		return errors.Wrap(err, "can't marshal a record")
	}
	if b.save(b.makeID(dirPath), jRecord) != nil {
		return errors.Wrap(err, "can't create a record in the db")
	}
	return nil
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}

// create a bolt id
func (b *BoltDB) makeID(dirPath string) string {
	return DirID(dirPath)
}

// save a record
func (b *BoltDB) save(id string, record []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dirsBucket))
		err := b.Put([]byte(id), record)
		return err
	})
	return err
}
