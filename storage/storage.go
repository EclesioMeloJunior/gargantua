package storage

import (
	"bytes"
	"errors"
	"path/filepath"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
)

const CHIANDB_FILE = "gg.db"

var (
	GeneralBoltBucket    = []byte("bolt_general")
	ErrNotImplementedYet = errors.New("function not implemented yet")
)

type Storage struct {
	db *leveldb.Database
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
	concreteKey := concatBucketAndKey(GeneralBoltBucket, key)
	return s.db.Has(concreteKey)
}

func (s *Storage) GetFromBucket(bucket, key []byte) ([]byte, error) {
	concreteKey := concatBucketAndKey(bucket, key)
	return s.db.Get(concreteKey)
}

func (s *Storage) DeleteFromBucket(bucket, key []byte) (err error) {
	concreteKey := concatBucketAndKey(bucket, key)
	return s.db.Delete(concreteKey)
}

func (s *Storage) PutOnBucket(bucket, key, value []byte) (err error) {
	concreteKey := concatBucketAndKey(bucket, key)
	return s.db.Put(concreteKey, value)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) EthereumDB() ethdb.KeyValueStore {
	return s.db
}

func NewStorage(basepath string) (*Storage, error) {
	dbfiles := filepath.Join(basepath, CHIANDB_FILE)

	st := new(Storage)

	var err error
	st.db, err = leveldb.New(dbfiles, 1024*100, 2, "", false)
	if err != nil {
		return nil, errors.New("problems to create/open a new database")
	}

	return st, nil
}

// concatBucketAndKey will
func concatBucketAndKey(bk, key []byte) []byte {
	return bytes.Join([][]byte{bk, {'-'}, key}, []byte{})
}
