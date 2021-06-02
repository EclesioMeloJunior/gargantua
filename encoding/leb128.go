package encoding

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
