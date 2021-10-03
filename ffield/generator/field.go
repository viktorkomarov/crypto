package main

import (
	"fmt"
	"math"

	"github.com/viktorkomarov/crypto/bitset"
	"github.com/viktorkomarov/crypto/ffield"
)

// TODO::support other degree ()
type GF2 struct {
	degree      uint
	polynomials uint64
	fields      map[uint64]*bitset.Set
}

func NewGF2(degree uint) (GF2, error) {
	if degree != 8 {
		return GF2{}, fmt.Errorf("degree out of range")
	}

	fields := make(map[uint64]*bitset.Set)
	for i := uint64(0); i < uint64(math.Pow(2.0, float64(degree))); i++ {
		fields[i] = bitset.SetFromNum(i)
	}

	return GF2{
		polynomials: findPolynomials(),
		degree:      degree,
		fields:      fields,
	}, nil
}

func findPolynomials() uint64 {
	return 283 // x8 + x4 + x3 + x1 + 1
}

func (g GF2) tableGenerator(f func(a, b uint64) uint64) map[ffield.Pair]uint64 {
	result := make(map[ffield.Pair]uint64)

	for a := range g.fields {
		for b := range g.fields {
			pair := ffield.NewPair(a, b)
			if _, ok := result[pair]; !ok {
				result[pair] = f(a, b)
			}
		}
	}

	return result
}

func (g GF2) generateSumTable() map[ffield.Pair]uint64 {
	return g.tableGenerator(func(a, b uint64) uint64 {
		return a ^ b
	})
}

func (g GF2) generateMulTable() map[ffield.Pair]uint64 {
	return g.tableGenerator(func(a, b uint64) uint64 {
		aSet, bSet := g.fields[a], g.fields[b]
		aSet = aSet.Mul(bSet)

		result := bitset.SetFromSize(int(g.degree))
		for _, n := range aSet.IndexOfOne() {
			temp := bitset.SetFromSize(n + 1)
			temp.SetVal(n, 1)

			if n >= int(g.degree) {
				temp = g.decreaseOrder(temp)
			}

			result = result.XOR(temp)
		}

		return result.BuildUint64()
	})
}

func (g GF2) decreaseOrder(set *bitset.Set) *bitset.Set {
	idx := uint(set.IndexOfOne()[0])
	if g.degree > idx {
		return set
	}

	diff := idx - g.degree
	diffSet := bitset.SetFromSize(int(g.degree))
	diffSet.SetVal(int(diff), 1)

	pSet := bitset.SetFromNum(g.polynomials)
	pSet.SetVal(int(g.degree), 0)

	return g.decreaseOrder(diffSet.Mul(pSet))
}
