package rlp_test

import (
	"testing"

	"github.com/EclesioMeloJunior/gargantua/rlp"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipisicing elit"
	enc := rlp.NewEncoder()
	_, err := enc.Encode(s)
	require.NoError(t, err)

	var dec string
	err = rlp.DecodeBytes(enc.Bytes(), &dec)
	require.NoError(t, err)

	require.Equal(t, s, dec)
}
