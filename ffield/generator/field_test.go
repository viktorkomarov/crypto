package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viktorkomarov/crypto/bitset"
)

func TestDecreaseOrder(t *testing.T) {
	testCases := []struct {
		desc     string
		expected func() *bitset.Set
		degree   func() *bitset.Set
	}{
		{
			desc: "x8",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(4, 1)
				result.SetVal(3, 1)
				result.SetVal(1, 1)
				result.SetVal(0, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(8, 1)
				return result
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			g := GF2{polynomials: findPolynomials(), degree: 8}
			actual := g.decreaseOrder(tC.degree())
			expected := tC.expected()
			for i := 0; i <= 7; i++ {
				require.Equalf(t, expected.Nth(i), actual.Nth(i), "at %d", i)
			}
		})
	}
}
