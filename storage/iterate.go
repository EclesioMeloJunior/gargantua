package storage

import (
	"github.com/boltdb/bolt"
)

type StorageIterator struct {
	currIdx int
	data    [][]byte

	err    error
	cursor *bolt.Cursor

	currKey   []byte
	currValue []byte
}

func (s *StorageIterator) Next() bool {
	if s.currIdx >= len(s.data)-1 {
		return true
	}

	s.currKey, s.currValue = s.data[s.currIdx], s.data[s.currIdx+1]
	s.currIdx += 2

	return false
}

func (s *StorageIterator) Key() []byte {
	return s.currKey
}

func (s *StorageIterator) Value() []byte {
	return s.currValue
}

func (s *StorageIterator) Error() error {
	return nil
}

func (s *StorageIterator) Release() {}

func (s *StorageIterator) add(k, v []byte) {
	s.data = append(s.data, k, v)
}
