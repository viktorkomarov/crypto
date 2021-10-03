package aes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viktorkomarov/crypto/ffield"
)

func TestAESTableMul(t *testing.T) {
	testCases := []struct {
		desc     string
		a        uint64
		b        uint64
		expected uint64
	}{
		{
			desc:     "157*56",
			a:        157,
			b:        56,
			expected: 177,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expected, ffield.MulGF8(tC.a, tC.b))
		})
	}
}

func TestAESTableSum(t *testing.T) {
	testCases := []struct {
		desc     string
		a        uint64
		b        uint64
		expected uint64
	}{
		{
			desc:     "157+56",
			a:        157,
			b:        56,
			expected: 165,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expected, ffield.SumGF8(tC.a, tC.b))
		})
	}
}
