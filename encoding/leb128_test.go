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
		{
			value:  624485,
			expect: []byte{0xE5, 0x8E, 0x26},
		},
		{
			value:  10,
			expect: []byte{0x0a},
		},
		{
			value:  90000,
			expect: []byte{0x90, 0xbf, 0x05},
		},
	}

	for _, test := range tests {
		out := encoding.EncodeUint64LEB128(test.value)
		assert.Equal(t, test.expect, out)
	}
}

func TestUintLebDecoding(t *testing.T) {
	tests := []struct {
		value  []byte
		expect uint64
	}{
		{
			value:  []byte{192, 62},
			expect: 8000,
		},
		{
			value:  []byte{0xE5, 0x8E, 0x26},
			expect: 624485,
		},
		{
			value:  []byte{0x0a},
			expect: 10,
		},
		{
			value:  []byte{0x90, 0xbf, 0x05},
			expect: 90000,
		},
	}

	for _, test := range tests {
		out, err := encoding.DecodeUint64LEB128(test.value)
		assert.Nil(t, err)
		assert.Equal(t, test.expect, out)
	}
}
