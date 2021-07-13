package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// DefaultKeystoreFile is a pattern to keys files name
const DefaultKeystoreFile = "%s/%s-%s.keystore"

func StoreKeyPair(name string, path string, pair *Pair, passowrd string) error {
	x509encodedPrivate, err := x509.MarshalECPrivateKey(pair.Private.PrivateKey)
	if err != nil {
		return err
	}

	x509encodedPublic, err := x509.MarshalPKIXPublicKey(pair.Public.PublicKey)
	if err != nil {
		return err
	}

	privateBuff := &bytes.Buffer{}
	err = pem.Encode(privateBuff, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509encodedPrivate,
	})
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	encryptedPrivateKey, err := aesEncrypt(privateBuff.Bytes(), passowrd)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	publicBuff := &bytes.Buffer{}
	err = pem.Encode(publicBuff, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509encodedPublic,
	})
	if err != nil {
		return fmt.Errorf("failed to encode public key: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf(DefaultKeystoreFile, path, name, PrivateType), encryptedPrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf(DefaultKeystoreFile, path, name, PublicType), publicBuff.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func aesEncrypt(data []byte, password string) ([]byte, error) {
	b, err := aes.NewCipher(hashpass(password))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func hashpass(password string) []byte {
	c := md5.New()
	c.Write([]byte(password))
	return c.Sum(nil)
}
