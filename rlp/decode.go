package rlp

import (
	"bytes"
	"errors"
	"io"
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
	ErrKindNotDefined = errors.New("rlp: kind could not be defined")

	decoderPool = sync.Pool{
		New: func() interface{} { return new(Decoder) },
	}
)

type Decoder struct {
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
		d.size = 0
		d.kind = Byte

	case b <= (0x80 + 55):
		d.size = uint64(b - 0x80)
		d.kind = String

	default:
		return ErrKindNotDefined
	}

	return nil
}

func DecodeBytes(b []byte, val interface{}) error {
	decoder := decoderPool.Get().(*Decoder)
	defer decoderPool.Put(decoder)

	decoder.r = bytes.NewReader(b)
	if err := decoder.defineKind(); err != nil {
		return err
	}

	return nil
}
