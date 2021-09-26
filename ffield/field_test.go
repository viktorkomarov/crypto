package ffield

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viktorkomarov/crypto/bitset"
)

func TestMul(t *testing.T) {
	testCases := []struct {
		desc     string
		a        *bitset.Set
		b        *bitset.Set
		expected *bitset.Set
	}{
		{
			desc:     "(x^3+x+1)*(x^2+x+1)=x^5+x^4+1",
			a:        bitset.SetFromBytes([]byte{0b11010000}),
			b:        bitset.SetFromBytes([]byte{0b11100000}),
			expected: bitset.SetFromBytes([]byte{0b10001100}),
		},
		{
			desc:     "(x^2+x+1)*(x^3+x+1)=x^5+x^4+1",
			b:        bitset.SetFromBytes([]byte{0b11010000}),
			a:        bitset.SetFromBytes([]byte{0b11100000}),
			expected: bitset.SetFromBytes([]byte{0b10001100}),
		},
		{
			desc:     "(x^5+x^3+x^2+1)*1=(x^5+x^3+x^2+1)",
			a:        bitset.SetFromBytes([]byte{0b10110100}),
			b:        bitset.SetFromBytes([]byte{0b10000000}),
			expected: bitset.SetFromBytes([]byte{0b10110100}),
		},
		{
			desc:     "(x^2+1)*(x^2+1)=(x^4+1)",
			a:        bitset.SetFromBytes([]byte{0b10100000}),
			b:        bitset.SetFromBytes([]byte{0b10100000}),
			expected: bitset.SetFromBytes([]byte{0b10001000}),
		},
		{
			desc:     "(x^3+x^2+1)*(x^2+x)=x^5+x^3+x^2+x",
			a:        bitset.SetFromBytes([]byte{0b10110000}),
			b:        bitset.SetFromBytes([]byte{0b01100000}),
			expected: bitset.SetFromBytes([]byte{0b01110100}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mul := Mul(tC.a, tC.b)
			for i := 0; i < tC.expected.Size(); i++ {
				require.Equalf(t, tC.expected.Nth(i), mul.Nth(i), "position %d", i)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	testCases := []struct {
		desc   string
		degree int
		count  int
	}{
		{
			desc:   "it works",
			degree: 8,
			count:  256,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.count, len(GenerateGF2(8)))
		})
	}
}
