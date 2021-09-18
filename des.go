package des

type Cipher struct {
	key *Bitset
}

func NewCipher(keySet []byte) Cipher {
	bitset := BitsetFromBytes(keySet)
	key := BitsetFromSize(56)
	for i := 0; i < 56; i++ {
		key.SetVal(i, bitset.Nth(keyInitTable[i]-1))
	}

	return Cipher{
		key: key,
	}
}

func (e Cipher) Encrypt(b []byte) []byte {
	return e.do(b, 0, 16)
}

func (e Cipher) Decrypt(msg []byte) []byte {
	return e.do(msg, 15, -1)
}

func (e Cipher) do(payload []byte, from, to int) []byte {
	msg := BitsetFromBytes(payload)
	bits := BitsetFromSize(64)
	for i := range initPermutationTable {
		bits.SetVal(i, msg.Nth(initPermutationTable[i]-1))
	}

	l, r := bits.Subset(0, 32), bits.Subset(32, 64)
	for i := from; i != to; {
		l, r = r, l.XOR(e.f(r, e.round(i)))
		if from < to {
			i++
		} else {
			i--
		}
	}
	l, r = r, l

	result := BitsetFromSize(64)
	for i := range initPermutatuinReverseTable {
		from := initPermutatuinReverseTable[i] - 1
		if from < 32 {
			result.SetVal(i, l.Nth(from))
		} else {
			result.SetVal(i, r.Nth(from-32))
		}
	}

	return result.Bits()
}

func buildByte(ns ...byte) byte {
	if len(ns) > 7 {
		return 0
	}

	var b byte
	for i := 0; i < len(ns); i++ {
		b |= (ns[i] << (len(ns) - i - 1))
	}
	return b
}

func (e Cipher) f(r, ki *Bitset) *Bitset {
	extendedR := BitsetFromSize(48)
	for i := range eBitSelection {
		extendedR.SetVal(i, r.Nth(eBitSelection[i]-1))
	}
	ki = ki.XOR(extendedR)

	sBoxed := BitsetFromSize(32)
	for i := 0; i < 8; i++ {
		s := ki.Subset(i*6, i*6+6)
		l := buildByte(s.Nth(0), s.Nth(5))

		bits := make([]byte, 4)
		for i := range bits {
			bits[i] = s.Nth(i + 1)
		}
		c := buildByte(bits...)

		val := byte(sTables[i][l][c])
		for j := 0; j < 4; j++ {
			sBoxed.SetVal(i*4+j, (val>>(3-j))&1)
		}
	}

	result := BitsetFromSize(32)
	for i := range pBitMutation {
		result.SetVal(i, sBoxed.Nth(pBitMutation[i]-1))
	}
	return result
}

func (e Cipher) round(r int) *Bitset {
	shiftOne := map[int]bool{0: true, 1: true, 8: true, 15: true}
	shift := 0
	for i := 0; i <= r; i++ {
		if shiftOne[i] {
			shift += 1
		} else {
			shift += 2
		}
	}

	c := e.key.Subset(0, 28).LeftRotate(shift)
	d := e.key.Subset(28, 56).LeftRotate(shift)

	result := BitsetFromSize(48)
	for i := range keyRoundTable {
		from := keyRoundTable[i] - 1
		if from < 28 {
			result.SetVal(i, c.Nth(from))
		} else {
			result.SetVal(i, d.Nth(from-28))
		}
	}

	return result
}
