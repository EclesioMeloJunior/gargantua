package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrKeyNotFound = errors.New("cannot found key")
)

// LoadPublicKey parse the encoded ecdsa private key file to in memory ecdsa.PrivateKey
func LoadPrivateKey(basepath, name string) (*PrivateKey, error) {
	keysdir := filepath.Join(basepath, "keys")
	privateKeyPath := fmt.Sprintf(DefaultKeystoreFile, keysdir, name, PrivateType)

	privateExists, err := checkFileStat(privateKeyPath)
	if err != nil {
		return nil, err
	}

	if !privateExists {
		return nil, ErrKeyNotFound
	}

	privateKeyData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(privateKeyData)
	ecdsaPrivateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		ecdsaPrivateKey,
	}, nil
}

// LoadPublicKey parse the encoded ecdsa public key file to in memory ecdsa.PublicKey
func LoadPublicKey(basepath, name string) (*PublicKey, error) {
	keysdir := filepath.Join(basepath, "keys")
	publicKeyPath := fmt.Sprintf(DefaultKeystoreFile, keysdir, name, PublicType)

	pubExists, err := checkFileStat(publicKeyPath)
	if err != nil {
		return nil, err
	}

	if !pubExists {
		return nil, ErrKeyNotFound
	}

	pubBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(pubBytes)
	genericPubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &PublicKey{
		genericPubKey.(*ecdsa.PublicKey),
	}, nil
}

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

func aesDecrypt(data []byte, password string) ([]byte, error) {
	c, err := aes.NewCipher(hashpass(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("gcm.NonceSize() greater then len(to_decrypt)")
	}

	nonce, chipertext := data[:nonceSize], data[nonceSize:]
	plaintxt, err := gcm.Open(nil, nonce, chipertext, nil)
	if err != nil {
		return nil, err
	}

	return []byte(plaintxt), nil
}
