package storage

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/ethdb"
)

type StorageBatcher struct {
	st   *Storage
	data map[string][]byte
}

func (s *StorageBatcher) Reset()                              {}
func (s *StorageBatcher) Replay(w ethdb.KeyValueWriter) error { return ErrNotImplementedYet }

func (s *StorageBatcher) Write() error {
	for k, v := range s.data {
		originalBytes, err := hex.DecodeString(k)
		if err != nil {
			return err
		}

		if err = s.st.Put(originalBytes, v); err != nil {
			return err
		}
	}

	s.data = map[string][]byte{}
	return nil
}

func (s *StorageBatcher) ValueSize() int {
	return len(s.data)
}

func (s *StorageBatcher) Put(key, value []byte) error {
	s.data[hex.EncodeToString(key)] = value
	return nil
}

func (s *StorageBatcher) Delete(key []byte) error {
	_, ok := s.data[hex.EncodeToString(key)]
	if ok {
		delete(s.data, hex.EncodeToString(key))
	}

	return nil
}
