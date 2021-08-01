package storage

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	GeneralBoltBucket    = []byte("bolt_general")
	ErrNotImplementedYet = errors.New("function not implemented yet")
)

type Storage struct {
	db *bolt.DB
}

func (s *Storage) Compact(start []byte, limit []byte) error { return ErrNotImplementedYet }

func (s *Storage) Stat(property string) (string, error) { return "", ErrNotImplementedYet }

func (s *Storage) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	boltTx, err := s.db.Begin(false)
	if err != nil {
		log.Println("problems to open a transaction while creating iterator")
		return nil
	}

	defer boltTx.Commit()

	it := &StorageIterator{
		cursor: nil,
		data:   make([][]byte, 0),
	}

	cursor := boltTx.Bucket(GeneralBoltBucket).Cursor()

	if prefix != nil || len(prefix) > 0 {
		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			it.add(k, v)
		}
	}

	return it
}

func (s *Storage) NewBatch() ethdb.Batch {
	return &StorageBatcher{
		data: make(map[string][]byte),
		st:   s,
	}
}

func (s *Storage) Put(key, value []byte) (err error) {
	return s.PutOnBucket(GeneralBoltBucket, key, value)
}

func (s *Storage) Delete(key []byte) (err error) {
	return s.DeleteFromBucket(GeneralBoltBucket, key)
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	return s.GetFromBucket(GeneralBoltBucket, key)
}

func (s *Storage) Has(key []byte) (bool, error) {
	v, err := s.GetFromBucket(GeneralBoltBucket, key)
	if err != nil {
		return false, err
	}

	return (v != nil) || len(v) > 0, nil
}

func (s *Storage) GetFromBucket(bucket, key []byte) ([]byte, error) {
	boltTx, err := s.db.Begin(false)
	if err != nil {
		return nil, err
	}

	defer boltTx.Commit()

	boltBucket := boltTx.Bucket(bucket)
	if boltBucket == nil {
		return nil, bolt.ErrBucketNotFound
	}

	return boltBucket.Get(key), nil
}

func (s *Storage) DeleteFromBucket(bucket, key []byte) (err error) {
	boltTx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := boltTx.Rollback(); rollbackErr != nil {
				err = rollbackErr
				return
			}

			return
		}

		if commitErr := boltTx.Commit(); commitErr != nil {
			err = commitErr
			return
		}
	}()

	boltBucket := boltTx.Bucket(bucket)
	if boltBucket == nil {
		return bolt.ErrBucketNotFound
	}

	return boltBucket.Delete(key)
}

func (s *Storage) PutOnBucket(bucket, key, value []byte) (err error) {
	boltTx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := boltTx.Rollback(); rollbackErr != nil {
				err = rollbackErr
				return
			}
			return
		}

		if commitErr := boltTx.Commit(); commitErr != nil {
			err = commitErr
			return
		}
	}()

	b, err := boltTx.CreateBucketIfNotExists(bucket)
	if err != nil {
		return err
	}

	if err := b.Put(key, value); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close() error {
	log.Println("beign called ... (we might have problems)")
	return s.db.Close()
}

func NewStorage(basepath string) (*Storage, error) {
	dbfiles := filepath.Join(basepath, "storage.db")
	_, err := os.Stat(dbfiles)

	if errors.Is(os.ErrNotExist, err) {
		return nil, errors.New("database files alreaady exists, choose another location to avoid overwriten")
	}

	st := new(Storage)
	st.db, err = bolt.Open(dbfiles, os.ModePerm, nil)
	if err != nil {
		return nil, errors.New("problems to open a new database")
	}

	return st, nil
}
