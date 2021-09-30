package main

import (
	"math"

	"github.com/viktorkomarov/crypto/bitset"
)

type pair struct {
	left  string
	right string
}

const gf2_2 = 0b11100000

func newPair(l, r *bitset.Set) pair {
	if r.String() < l.String() {
		l, r = r, l
	}

	return pair{
		left:  l.String(),
		right: r.String(),
	}
}

func (p pair) String() string {
	return p.left + "_" + p.right
}

type GF2 struct {
	degree int
	fields map[string]*bitset.Set
}

func NewGF2(degree int) GF2 {
	return GF2{
		degree: degree,
		fields: generateGF2(degree),
	}
}

// change to copy
func (g *GF2) Field() map[string]*bitset.Set {
	return g.fields
}

func (g *GF2) GenerateSumTable(field map[string]*bitset.Set) map[string]*bitset.Set {
	sumTable := make(map[string]*bitset.Set)

	for _, lSet := range field {
		for _, rSet := range field {
			pr := newPair(lSet, rSet)
			if _, ok := sumTable[pr.String()]; !ok {
				sumTable[pr.String()] = lSet.XOR(rSet)
			}
		}
	}

	return sumTable
}

func (g *GF2) GenerateMulTable(field map[string]*bitset.Set) map[string]*bitset.Set {
	mulTable := make(map[string]*bitset.Set)

	for _, lSet := range field {
		for _, rSet := range field {
			pr := newPair(lSet, rSet)
			if _, ok := mulTable[pr.String()]; !ok {
				mulTable[pr.String()] = g.DivRem(g.Mul(lSet, rSet), bitset.SetFromBytes([]byte{gf2_2}))
			}
		}
	}

	return mulTable
}

func generateGF2(degree int) map[string]*bitset.Set {
	countNums := int(math.Pow(2.0, float64(degree)))
	result := make(map[string]*bitset.Set, countNums)

	for i := 0; i < countNums; i++ {
		set := bitset.SetFromSize(degree)

		for j := 0; j < 8; j++ {
			set.SetVal(j, byte((i>>j)&1))
		}

		result[set.String()] = set
	}

	return result
}

func (g *GF2) Mul(a, b *bitset.Set) *bitset.Set {
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

func (g *GF2) DivRem(dividend, divider *bitset.Set) *bitset.Set {
	div1, div2 := dividend.IndexOfOne(), dividend.IndexOfOne()
	for _, b := range []bool{len(div1) == 0, len(div2) == 0, div1[0] < div2[0]} {
		if b {
			return nil
		}
	}

	return g.decreaseDegree(div1[0], dividend, divider)
}

// a / b , maxi - index one set of a
func (g *GF2) decreaseDegree(maxi int, a, b *bitset.Set) *bitset.Set {
	if maxi < g.degree {
		return a
	}

	diff := bitset.SetFromSize(g.degree)
	diff.SetVal(maxi-g.degree, 1)
	diff = g.Mul(diff, b)
	diff.SetVal(maxi, 0)

	result := bitset.SetFromSize(g.degree)
	indexOfOne := diff.IndexOfOne()
	for _, idx := range indexOfOne {
		setter := bitset.SetFromSize(g.degree)
		setter.SetVal(idx, 1)
		result = result.XOR(g.decreaseDegree(idx, setter, b))
	}

	return result
}
