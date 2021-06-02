package encoding

import "errors"

// EncodeUint64LEB128 compress a unsigned interger in a small number of bytes
func EncodeUint64LEB128(v uint64) []byte {
	b := []byte{}

	for {
		// spliting in 7 bit group
		group := uint8(v & 0x7f) // (hex 0x7f) (binary 1111111) (decimal 127)

		// move to the next group
		v >>= 7

		if v != 0 {
			group |= 0x80 // high one in the if there is more groups to evaluate
		}

		b = append(b, group)

		if group&0x80 == 0 {
			break
		}
	}

	return b
}

func DecodeUint64LEB128(b []byte) (uint64, error) {
	l := uint8(len(b) & 0xff)
	if l > 10 {
		return 0, errors.New("the max leb128 encoded bytes len is 10")
	}

	var result uint64
	for i := uint8(0); i < l; i++ {
		result |= uint64(b[i]&0x7f) << (7 * i)
		if b[i]&0x80 == 0 {
			break
		}
	}

	return result, nil
}
