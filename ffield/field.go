package ffield

import "github.com/viktorkomarov/crypto/bitset"

func Mul(a, b *bitset.Set) *bitset.Set {
	aOnes, bOnes := a.IndexOfOne(), b.IndexOfOne()
	if len(aOnes) == 0 || len(bOnes) == 0 {
		return bitset.SetFromSize(0)
	}

	mul := bitset.SetFromSize(aOnes[0] + bOnes[0])
	for _, ai := range aOnes {
		for _, bi := range bOnes {
			mul.SetVal(ai+bi, mul.Nth(ai+bi)^1)
		}
	}

	return mul
}
