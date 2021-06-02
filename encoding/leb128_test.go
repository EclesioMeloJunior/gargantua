package encoding_test

import (
	"testing"

	"github.com/EclesioMeloJunior/gargantua/encoding"
	"github.com/stretchr/testify/assert"
)

func TestUintLebEncoding(t *testing.T) {
	tests := []struct {
		value  uint64
		expect []byte
	}{
		{
			value:  8000,
			expect: []byte{192, 62},
		},
	}

	for _, test := range tests {
		out := encoding.EncodeUint64LEB128(test.value)
		assert.Equal(t, test.expect, out)
	}
}
