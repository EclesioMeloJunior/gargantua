package rlp

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

var (
	EmptyString = []byte{0x80}
	EmptyList   = []byte{0xc0}
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
		return e.encodeBytes(t, EmptyString[0])
	case []int, []uint:
		return e.encodeNumber(t)
	case []string:
		enc := NewEncoder()
		for _, s := range t {
			_, err := enc.Encode(s)
			// TODO: return the total of bytes written
			if err != nil {
				return 0, err
			}
		}
		return e.encodeBytes(enc.Bytes(), EmptyList[0])
	case [][]byte:
		enc := NewEncoder()
		for _, s := range t {
			fmt.Println(s)
			_, err := enc.Encode(s)
			// TODO: return the total of bytes written
			if err != nil {
				return 0, err
			}
			fmt.Println(enc.Bytes())
		}
		return e.encodeBytes(enc.Bytes(), EmptyList[0])
	default:
		rval := reflect.ValueOf(b)
		writer, err := cachedWriter(rval.Type())
		if err != nil {
			return 0, fmt.Errorf("unsuported %s type to encode", t)
		}

		return writer(rval, e)
	}
}

func (e *Encoder) encodeBytes(i interface{}, offset byte) (int, error) {
	d, err := fromStringToBytes(i)
	if err != nil {
		return 0, err
	}

	// if there is just one item and
	// this byte is in the range [0x00, 0x7f]
	if len(d) == 1 && (d[0]&offset) == 0 {
		return e.buff.Write(d)
	}

	// if b is a 0-55 len bytes long,
	if len(d) < 56 {
		first := offset + byte(len(d))
		all := bytes.Join([][]byte{{first}, d}, []byte{})
		return e.buff.Write(all)
	}

	// if b is greater than 55 bytes
	base2 := binaryForm(len(d))
	first := byte(len(base2)) + offset + 55

	all := bytes.Join([][]byte{{first}, {byte(len(d))}, d}, []byte{})
	return e.buff.Write(all)
}

func (e *Encoder) encodeNumber(i interface{}) (n int, err error) {
	return 0, nil
}

func (e *Encoder) Bytes() []byte {
	return e.buff.Bytes()
}

func fromStringToBytes(i interface{}) (b []byte, err error) {
	switch i := i.(type) {
	case string:
		b = []byte(i)
	case []byte:
		b = i
	default:
		err = errors.New("argument must be string or byte array")
	}

	return
}

func binaryForm(i int) []byte {
	if i == 0 {
		return []byte{}
	}

	return bytes.Join([][]byte{binaryForm(i / 256), {byte(i % 256)}}, []byte{})
}

func createWriter(t reflect.Type) (writer, error) {
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		return createPtrWriter(t)
	case reflect.Struct:
		return
	default:
		return nil, fmt.Errorf("writer not implemented for type %s", t)
	}
}

func createPtrWriter(t reflect.Type) (writer, error) {
	info := new(typeinfo)
	info.generate(t.Elem())

	if info.writerErr != nil {
		return nil, info.writerErr
	}

	writer := func(val reflect.Value, e *Encoder) (int, error) {
		if val.IsNil() {
			return 0, errors.New("could not encode nil pointer")
		}

		return info.writer(val.Elem(), e)
	}

	return writer, nil
}

func createStructWriter(t reflect.Type) (writer, error) {
	fields, err := getStructFields(t)
}
