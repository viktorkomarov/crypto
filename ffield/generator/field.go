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
	mulTable    map[ffield.Pair]uint64
	fields      map[uint64]*bitset.Set
}

func NewGF2(degree uint) (*GF2, error) {
	if degree != 8 {
		return nil, fmt.Errorf("degree out of range")
	}

	fields := make(map[uint64]*bitset.Set)
	for i := uint64(0); i < uint64(math.Pow(2.0, float64(degree))); i++ {
		fields[i] = bitset.SetFromNum(i)
	}

	return &GF2{
		polynomials: findPolynomials(),
		degree:      degree,
		fields:      fields,
	}, nil
}

func findPolynomials() uint64 {
	return 283 // x8 + x4 + x3 + x1 + 1
}

func (g *GF2) tableGenerator(f func(a, b uint64) uint64) map[ffield.Pair]uint64 {
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

func (g *GF2) generateSumTable() map[ffield.Pair]uint64 {
	return g.tableGenerator(func(a, b uint64) uint64 {
		return a ^ b
	})
}

func (g *GF2) generateMulTable() map[ffield.Pair]uint64 {
	if g.mulTable == nil {
		g.mulTable = g.tableGenerator(func(a, b uint64) uint64 {
			aSet, bSet := g.fields[a], g.fields[b]
			return g.decreaseOrderIfNeed(aSet.Mul(bSet)).BuildUint64()
		})
	}

	return g.mulTable
}

func (g *GF2) decreaseOrderIfNeed(set *bitset.Set) *bitset.Set {
	result := bitset.SetFromSize(int(g.degree))

	degress := set.IndexOfOne()
	for len(degress) > 0 {
		idx := degress[0]
		if g.degree <= uint(idx) {
			diff := uint(idx) - g.degree
			diffSet := bitset.SetFromSize(int(g.degree))
			diffSet.SetVal(int(diff), 1)
			pSet := bitset.SetFromNum(g.polynomials)
			pSet.SetVal(int(g.degree), 0)
			degress = append(degress, diffSet.Mul(pSet).IndexOfOne()...)
		} else {
			temp := bitset.SetFromSize(idx + 1)
			temp.SetVal(idx, 1)
			result = result.XOR(temp)
		}
		degress = degress[1:]
	}

	return result
}

func (g *GF2) generateInvrTable() map[uint64]uint64 {
	mul := g.mulTable
	if mul == nil {
		mul = g.generateMulTable()
	}

	result := make(map[uint64]uint64)
	for k, v := range mul {
		if v == 1 {
			result[k.L] = k.R
			result[k.R] = k.L
		}
	}

	return result
}
