package ffield

import (
	"math"

	"github.com/viktorkomarov/crypto/bitset"
)

type pair struct {
	left  string
	right string
}

func newPair(l, r *bitset.Set) pair {
	if r.String() < l.String() {
		l, r = r, l
	}

	return pair{
		left:  l.String(),
		right: r.String(),
	}
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

func (g *GF2) GenerateSumTable(field map[string]*bitset.Set) map[pair]*bitset.Set {
	sumTable := make(map[pair]*bitset.Set)

	for _, lSet := range field {
		for _, rSet := range field {
			pr := newPair(lSet, rSet)
			if _, ok := sumTable[pr]; !ok {
				sumTable[pr] = lSet.XOR(rSet)
			}
		}
	}

	return sumTable
}

func (g *GF2) GenerateMulTable(field map[string]*bitset.Set) map[pair]*bitset.Set {
	mulTable := make(map[pair]*bitset.Set)

	for _, lSet := range field {
		for _, rSet := range field {
			pr := newPair(lSet, rSet)
			if _, ok := mulTable[pr]; !ok {
				mulTable[pr] = Div(Mul(lSet, rSet), nil)
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

	rem := bitset.SetFromSize(g.degree)
	for div1[0] > div2[0] {
		diff := div1[0] - div2[0]
		rem.SetVal(diff, 1)
		temp := bitset.SetFromSize(g.degree)
		temp.SetVal(diff, 1)
		div1 = g.Mul(divider, temp).IndexOfOne()
	}

	return rem
}
