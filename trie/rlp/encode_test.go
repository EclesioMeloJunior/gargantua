package rlp_test

import (
	"testing"

	"github.com/EclesioMeloJunior/gargantua/trie/rlp"
	"github.com/stretchr/testify/require"
)

var tests = []struct {
	s string
	e []byte
}{
	{
		s: "",
		e: []byte{0x80},
	},
	{
		s: "dog",
		e: []byte{byte(0x83), byte('d'), byte('o'), byte('g')},
	},
}

func TestRLPEncoding(t *testing.T) {
	for _, s := range tests {
		e := rlp.NewEncoder()
		_, err := e.Encode(s.s)

		require.NoError(t, err)
		require.Equal(t, e.Bytes(), s.e)
	}
}

func TestRLPEncodingMore55BytesLen(t *testing.T) {
	b := make([]byte, 56)
	for i := 0; i < 56; i++ {
		b[i] = byte(i)
	}

	e := rlp.NewEncoder()
	_, err := e.Encode(b)
	require.NoError(t, err)

	exp := []byte{byte(0xbd), byte(56)}
	exp = append(exp, b...)

	require.Equal(t, exp, e.Bytes())
}
