package keystore_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/stretchr/testify/require"
)

func TestNewKeyPair_PrivateKeySig_AndPublicKeyVerification(t *testing.T) {
	kp, err := keystore.NewKeyPair()
	require.NoError(t, err)

	message := []byte("My message to be signed")
	messagehash := sha256.Sum256(message)

	sig, err := ecdsa.SignASN1(rand.Reader, kp.Private.PrivateKey, messagehash[:])
	require.NoError(t, err)

	fmt.Printf("%x\n", sig)
	valid := ecdsa.VerifyASN1(&kp.Private.PublicKey, messagehash[:], sig)
	require.True(t, valid)
}
