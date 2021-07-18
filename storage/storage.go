package storage

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

type Storage struct {
	db *bolt.DB
}

func NewStorage(basepath string) (*Storage, error) {
	dbfiles := filepath.Join(basepath, "storage.db")
	_, err := os.Stat(dbfiles)

	if errors.Is(os.ErrNotExist, err) {
		return nil, errors.New("database files alreaady exists, choose another location to avoid overwriten")
	}

	var st *Storage

	st.db, err = bolt.Open(dbfiles, os.ModePerm, nil)
	if err != nil {
		return nil, errors.New("problems to open a new database")
	}

	return st, nil
}
