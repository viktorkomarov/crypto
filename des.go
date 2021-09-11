package des

type Encoder struct {
	key *Bitset
}

func NewEncoder(keySet []byte) Encoder {
	bitset := BitsetFromBytes(keySet)
	key := BitsetFromSize(56)
	for i := 0; i < 56; i++ {
		key.SetVal(i, bitset.Nth(keyInitTable[i]-1))
	}

	return Encoder{
		key: key,
	}
}

func (e Encoder) Encode(b []byte) []byte {
	msg := BitsetFromBytes(b)
	bits := BitsetFromSize(64)
	for i := range initPermutationTable {
		bits.SetVal(i, msg.Nth(initPermutationTable[i]-1))
	}

	l, r := bits.Subset(0, 32), bits.Subset(32, 64)
	for i := 0; i < 16; i++ {
		l, r = r, l.XOR(e.f(r, e.round(i)))
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
		b |= (ns[i] << i)
	}
	return b
}

func (e Encoder) f(r, ki *Bitset) *Bitset {
	extendedR := BitsetFromSize(48)
	for i := range eBitSelection {
		extendedR.SetVal(i, r.Nth(eBitSelection[i]-1))
	}
	xor := ki.XOR(extendedR)

	sboxs := make([]byte, 0, 8)
	for i := 0; i < 8; i++ {
		s := xor.Subset(i*6, i*6+6)
		l, c := buildByte(s.Nth(0), s.Nth(5)), buildByte(s.Subset(1, 5).Bits()...)
		sboxs = append(sboxs, byte(sTables[i][l][c]))
	}

	sBoxed := BitsetFromBytes(sboxs)
	result := BitsetFromSize(32)
	for i := range pBitMutation {
		result.SetVal(i, sBoxed.Nth(pBitMutation[i]-1))
	}
	return result
}

func (e Encoder) round(r int) *Bitset {
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
