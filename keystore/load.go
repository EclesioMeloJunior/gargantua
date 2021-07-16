package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CheckNodeHasKeyPair(basepath, name string) (bool, error) {
	keysdir := filepath.Join(basepath, "keys")

	publicKeyPath := fmt.Sprintf(DefaultKeystoreFile, keysdir, name, PublicType)
	privateKeyPath := fmt.Sprintf(DefaultKeystoreFile, keysdir, name, PrivateType)

	pubExists, err := checkFileStat(publicKeyPath)
	if err != nil {
		return false, err
	}

	privExists, err := checkFileStat(privateKeyPath)
	if err != nil {
		return false, err
	}

	return pubExists && privExists, nil
}

func checkFileStat(filepath string) (bool, error) {
	finfo, err := os.Stat(filepath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, errors.New("node doesnt have key pair")
	}

	return finfo != nil, nil
}
