package keystore_test

import (
	"testing"

	"github.com/EclesioMeloJunior/gargantua/keystore"
	"github.com/stretchr/testify/require"
)

func TestGetAddress_FromPublicKey(t *testing.T) {
	kp, err := keystore.NewKeyPair()
	require.NoError(t, err)

	addr := keystore.GetAddress(kp.Public)

	require.Len(t, addr.String()[2:], 40)
}
