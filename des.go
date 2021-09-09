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

func (e Encoder) round(r int) []byte {
	shiftOne := map[int]bool{1: true, 8: true, 15: true}
	shift := 1
	for i := 1; i < r; i++ {
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

	return result.Bits()
}
