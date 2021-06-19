package rlp

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

type Encoder struct {
	buff bytes.Buffer
}

func NewEncoder() *Encoder {
	e := new(Encoder)
	e.buff = bytes.Buffer{}

	return e
}

func (e *Encoder) Encode(b interface{}) (int, error) {
	switch t := b.(type) {
	case string, []byte:
		return e.encodeString(t)
	case []int, []uint, []string, [][]byte:
		return 0, fmt.Errorf("unimplemented %s type to encode", t)
	default:
		return 0, fmt.Errorf("unsuported %s type to encode", t)
	}
}

func (e *Encoder) encodeString(i interface{}) (n int, err error) {
	var d []byte

	switch i := i.(type) {
	case string:
		d = []byte(i)
	case []byte:
		d = i
	}

	// if there is just one item and
	// this byte is in the range [0x00, 0x7f]
	if len(d) == 1 && (d[0]&0x80) == 0 {
		return e.buff.Write(d)
	}

	// if b is a 0-55 len bytes long,
	if len(d) < 56 {
		first := byte(0x80 + len(d))
		all := bytes.Join([][]byte{{first}, d}, []byte{})
		return e.buff.Write(all)
	} else if big.NewInt(int64(len(d))).Cmp(new(big.Int).Exp(big.NewInt(256), big.NewInt(8), nil)) == -1 {
		base2 := strconv.FormatInt(int64(len(d)), 2)
		fmt.Println(base2, len(d))
		first := byte(0xb7 + len(base2))
		all := bytes.Join([][]byte{{first}, {byte(len(d))}, d}, []byte{})
		return e.buff.Write(all)
	}

	return 0, errors.New("too large input")
}

func (e *Encoder) Bytes() []byte {
	return e.buff.Bytes()
}
