package ffield

import "github.com/viktorkomarov/crypto/bitset"

func Mul(a, b *bitset.Set) *bitset.Set {
	result := bitset.SetFromSize(a.Size() + b.Size())
	for i := a.Size() - 1; i >= 0; i-- {
		if a.Nth(i) == 1 {
			for j := b.Size() - 1; j >= 0; j-- {
				if b.Nth(j) == 1 {
					result.SetVal(i+j, result.Nth(i+j)^1)
				}
			}
		}
	}

	return result
}
