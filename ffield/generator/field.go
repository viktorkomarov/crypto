package main

import (
	"fmt"
	"math"

	"github.com/viktorkomarov/crypto/bitset"
)

// TODO::support other degree ()
type GF2 struct {
	degree      uint
	polynomials uint64
	fields      map[uint64]*bitset.Set
}

type Pair struct {
	L, R uint64
}

func newPair(l, r uint64) Pair {
	if r < l {
		l, r = r, l
	}

	return Pair{
		L: l,
		R: r,
	}
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

func (g GF2) tableGenerator(f func(a, b uint64) uint64) map[Pair]uint64 {
	result := make(map[Pair]uint64)

	for a, _ := range g.fields {
		for b, _ := range g.fields {
			pair := newPair(a, b)
			if _, ok := result[pair]; !ok {
				result[pair] = f(a, b)
			}
		}
	}

	return result
}

func (g GF2) generateSumTable() map[Pair]uint64 {
	return g.tableGenerator(func(a, b uint64) uint64 {
		return a ^ b
	})
}

func (g GF2) generateMulTable() map[Pair]uint64 {
	return g.tableGenerator(func(a, b uint64) uint64 {
		return (a * b) % g.polynomials
	})
}
