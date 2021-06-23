package rlp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
)

type ByteReader interface {
	io.Reader
	io.ByteReader
}

type Kind uint8

const (
	Byte Kind = iota
	String
	List
)

var (
	ErrInvalidTypeSize = errors.New("rlp: type size is invalid")
	ErrNoPointer       = errors.New("rlp: decode to non pointer invalid")
	ErrDecodeIntoNil   = errors.New("rlp: cannot decode into nil interface")
	ErrGreaterByteSize = errors.New("rlp: greater byte size could no bte 0")

	decoderPool = sync.Pool{
		New: func() interface{} { return new(Decoder) },
	}
)

type Decoder struct {
	byteval byte

	size uint64
	kind Kind
	r    ByteReader
}

func (d *Decoder) defineKind() error {
	b, err := d.r.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case b < 0x80:
		d.byteval = b
		d.size = 0
		d.kind = Byte

	case b <= (0x80 + 55):
		d.size = uint64(b - 0x80)
		d.kind = String

	case b < 0xc0:
		size, err := d.checkGreaterSize(uint64(b - (0x80 + 55)))
		if err != nil {
			return err
		}
		d.size = size
		d.kind = String

	case b <= (0xc0 + 55):
		d.size = uint64(b - 0xc0)
		d.kind = List

	default:
		size, err := d.checkGreaterSize(uint64(b - (0xc0 + 55)))
		if err != nil {
			return err
		}

		d.size = size
		d.kind = List
	}

	return nil
}

func (d *Decoder) checkGreaterSize(s uint64) (uint64, error) {
	switch s {
	case 0:
		return 0, ErrGreaterByteSize
	case 1:
		lenbyte, err := d.r.ReadByte()
		return uint64(lenbyte), err
	default:
		buff := make([]byte, 8)
		for i := range buff {
			buff[i] = 0
		}

		start := int(8 - s)
		if err := d.readFull(buff[start:]); err != nil {
			return 0, err
		}

		if buff[start] == 0 {
			return 0, ErrGreaterByteSize
		}

		return binary.BigEndian.Uint64(buff[:]), nil
	}
}

func (d *Decoder) readFull(b []byte) error {
	var n, nn int
	var err error

	for n < len(b) && err == nil {
		nn, err = d.r.Read(b[n:])
		n += nn
	}

	return err
}

func (d *Decoder) Decode(val interface{}) error {
	if val == nil {
		return ErrDecodeIntoNil
	}

	rval := reflect.ValueOf(val)
	rtyp := rval.Type()
	if rtyp.Kind() != reflect.Ptr {
		return ErrNoPointer
	}

	if rval.IsNil() {
		return ErrDecodeIntoNil
	}

	decoder, err := cachedDecoder(rtyp.Elem())
	if err != nil {
		return err
	}

	return decoder(d, rval.Elem())
}

func (d *Decoder) Bytes() ([]byte, error) {
	switch d.kind {
	case Byte:
		return []byte{d.byteval}, nil
	case String:
		b := make([]byte, d.size)
		if err := d.readFull(b); err != nil {
			return nil, err
		}

		if d.size == 1 && b[0] < 0x80 {
			return nil, ErrInvalidTypeSize
		}

		return b, nil
	default:
		return nil, errors.New("expected string to decoder")
	}
}

func DecodeBytes(b []byte, val interface{}) error {
	decoder := decoderPool.Get().(*Decoder)
	defer decoderPool.Put(decoder)

	decoder.r = bytes.NewReader(b)
	if err := decoder.defineKind(); err != nil {
		return err
	}

	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}

func createDecoder(t reflect.Type) (decoder, error) {
	k := t.Kind()
	switch k {
	case reflect.String:
		return decodeString, nil
	default:
		return nil, fmt.Errorf("decoder not implemented for type %s", k)
	}
}

func decodeString(dec *Decoder, v reflect.Value) error {
	b, err := dec.Bytes()
	if err != nil {
		return err
	}

	v.SetString(string(b))
	return nil
}
