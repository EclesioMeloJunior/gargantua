package config

import (
	"os"
	"path/filepath"
	"strings"
)

//SetupBasepath creates the basepath dir if it not exists
func SetupBasepath(basepath string) error {
	dir, err := ExpandDir(basepath)
	if err != nil {
		return err
	}

	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	} else if !os.IsNotExist(err) && err != nil {
		return err
	}

	return createKeysDir(basepath)
}

func ExpandDir(dir string) (string, error) {
	if strings.HasPrefix(dir, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dir = strings.Replace(dir, "~", homedir, -1)
	}

	return dir, nil
}

func createKeysDir(basepaht string) error {
	keysdir := filepath.Join(basepaht, "keys")

	_, err := os.Stat(keysdir)

	if os.IsNotExist(err) {
		return os.MkdirAll(keysdir, os.ModeDir|os.ModePerm)
	}

	return err
}
