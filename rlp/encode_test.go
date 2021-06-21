package rlp_test

import (
	"testing"

	"github.com/EclesioMeloJunior/gargantua/rlp"
	"github.com/stretchr/testify/require"
)

func TestRLPSimpleEncoding(t *testing.T) {
	tests := []struct {
		s string
		e []byte
	}{
		{
			s: "",
			e: []byte{0x80},
		},
		{
			s: "dog",
			e: []byte{0x83, 'd', 'o', 'g'},
		},
	}

	for _, s := range tests {
		e := rlp.NewEncoder()
		_, err := e.Encode(s.s)

		require.NoError(t, err)
		require.Equal(t, e.Bytes(), s.e)
	}
}

func TestRLPEncodingMore55BytesLen(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipisicing elit"
	exp := []byte{0xb8, 0x38}
	exp = append(exp, []byte(s)...)

	e := rlp.NewEncoder()
	_, err := e.Encode(s)
	require.NoError(t, err)

	require.Equal(t, exp, e.Bytes())
}

func TestRLPEncodingSlices(t *testing.T) {
	t1 := [][]byte{{'d', 'o', 'g'}, {'c', 'a', 't'}}
	exp := []byte{byte(0xc0 + 8), byte(0x80 + 3), 'd', 'o', 'g', byte(0x80 + 3), 'c', 'a', 't'}

	enc := rlp.NewEncoder()
	_, err := enc.Encode(t1)
	require.NoError(t, err)

	require.Equal(t, exp, enc.Bytes())
}
