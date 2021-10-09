package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viktorkomarov/crypto/bitset"
	"github.com/viktorkomarov/crypto/ffield"
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
		{
			desc: "x9",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(5, 1)
				result.SetVal(4, 1)
				result.SetVal(2, 1)
				result.SetVal(1, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(9, 1)
				return result
			},
		},
		{
			desc: "x10",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(6, 1)
				result.SetVal(5, 1)
				result.SetVal(3, 1)
				result.SetVal(2, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(10, 1)
				return result
			},
		},
		{
			desc: "x6",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(6, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(6, 1)
				return result
			},
		},
		{
			desc: "x13",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(6, 1)
				result.SetVal(0, 1)
				result.SetVal(3, 1)
				result.SetVal(2, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(13, 1)
				return result
			},
		},
		{
			desc: "x12",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(7, 1)
				result.SetVal(5, 1)
				result.SetVal(3, 1)
				result.SetVal(1, 1)
				result.SetVal(0, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(12, 1)
				return result
			},
		},
		{
			desc: "x11",
			expected: func() *bitset.Set {
				result := bitset.SetFromSize(10)
				result.SetVal(7, 1)
				result.SetVal(6, 1)
				result.SetVal(4, 1)
				result.SetVal(3, 1)
				return result
			},
			degree: func() *bitset.Set {
				result := bitset.SetFromSize(5)
				result.SetVal(11, 1)
				return result
			},
		},
		{
			desc: "(x^7+x^6+x)*(x^5+x^3+x^2+x+1)",
			expected: func() *bitset.Set {
				return bitset.SetFromNum(1)
			},
			degree: func() *bitset.Set {
				return bitset.SetFromBytes([]byte{0b01111000, 0b00111000})
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			g := GF2{polynomials: findPolynomials(), degree: 8}
			actual := g.decreaseOrderIfNeed(tC.degree())
			expected := tC.expected()
			for i := 0; i <= 7; i++ {
				require.Equalf(t, expected.Nth(i), actual.Nth(i), "at %d", i)
			}
		})
	}
}

func TestTableMul(t *testing.T) {
	testCases := []struct {
		desc     string
		a        uint64
		b        uint64
		expected uint64
	}{
		{
			desc:     "194 && 47",
			a:        194,
			b:        47,
			expected: 1,
		},
		{
			desc:     "208 && 122",
			a:        208,
			b:        122,
			expected: 1,
		},
	}

	g, err := NewGF2(8)
	require.NoError(t, err)
	mul := g.generateMulTable()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expected, mul[ffield.NewPair(tC.a, tC.b)])
		})
	}
}

func Test(t *testing.T) {
	testCases := []struct {
		desc     string
		expected int
	}{
		{
			desc:     "size test",
			expected: 255,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			g, err := NewGF2(8)
			require.NoError(t, err)
			require.Equal(t, tC.expected, len(g.generateInvrTable()))
		})
	}
}
