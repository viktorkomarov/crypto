package des

type Encoder struct {
	key *Bitset
}

func NewEncoder(keySet []byte) Encoder {
	bitset := BitsetFromBytes(keySet)
	key := BitsetFromSize(56)
	for i := 0; i < 56; i++ {
		key.SetVal(i, bitset.Nth(keyInitTable[i]))
	}

	return Encoder{
		key: key,
	}
}

func (e Encoder) round(r int) []byte {
	shiftOne := map[int]bool{0: true, 1: true, 8: true, 15: true}
	shift := 1
	for i := 0; i < r; i++ {
		if shiftOne[i] {
			shift += 1
		} else {
			shift += 2
		}
	}

	c := e.key.Subset(0, 28).LeftRotate(shift)
	c.Append(e.key.Subset(28, 56).LeftRotate(shift))

	result := BitsetFromSize(48)
	for i := range keyRoundTable {
		result.SetVal(i, c.Nth(keyInitTable[i]))
	}
	return result.Bits()
}
