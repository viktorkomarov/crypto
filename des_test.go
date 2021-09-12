package des

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyEncoder(t *testing.T) {
	testCases := []struct {
		desc     string
		key      []byte
		expected []byte
	}{
		{
			desc: "it works#1",
			key: []byte{
				0b00010011,
				0b00110100,
				0b01010111,
				0b01111001,
				0b10011011,
				0b10111100,
				0b11011111,
				0b11110001,
			},
			expected: []byte{
				0b11110000,
				0b11001100,
				0b10101010,
				0b11110101,
				0b01010110,
				0b01100111,
				0b10001111,
			},
		},
		{
			desc: "it works#2",
			key: []byte{
				0b00010011,
				0b00110100,
				0b01010111,
				0b01111001,
				0b10011011,
				0b10111100,
				0b11011111,
				0b11110001,
			},
			expected: []byte{
				0b11110000,
				0b11001100,
				0b10101010,
				0b11110101,
				0b01010110,
				0b01100111,
				0b10001111,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expected, NewEncoder(tC.key).key.Bits())
		})
	}
}

func TestKeyRounds(t *testing.T) {
	key := []byte{
		0b00010011,
		0b00110100,
		0b01010111,
		0b01111001,
		0b10011011,
		0b10111100,
		0b11011111,
		0b11110001,
	}

	testCases := []struct {
		desc     string
		round    int
		expected []byte
	}{
		{
			desc:  "round 1",
			round: 0,
			expected: []byte{
				0b00011011,
				0b00000010,
				0b11101111,
				0b11111100,
				0b01110000,
				0b01110010,
			},
		},
		{
			desc:  "last round",
			round: 15,
			expected: []byte{
				0b11001011,
				0b00111101,
				0b10001011,
				0b00001110,
				0b00010111,
				0b11110101,
			},
		},
		{
			desc:  "round 8",
			round: 7,
			expected: []byte{
				0b11110111,
				0b10001010,
				0b00111010,
				0b11000001,
				0b00111011,
				0b11111011,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := NewEncoder(key)
			require.Equal(t, tC.expected, bits.round(tC.round).Bits())
		})
	}
}

func TestBuildByte(t *testing.T) {
	testCases := []struct {
		desc     string
		args     []byte
		expected byte
	}{
		{
			desc:     "0 bytes",
			args:     []byte{},
			expected: 0,
		},
		{
			desc:     "1",
			args:     []byte{1},
			expected: 1,
		},
		{
			desc:     "10",
			args:     []byte{1, 0},
			expected: 2,
		},
		{
			desc:     "01",
			args:     []byte{0, 1},
			expected: 1,
		},
		{
			desc:     "1001",
			args:     []byte{1, 0, 0, 1},
			expected: 9,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expected, buildByte(tC.args...))
		})
	}
}

func TestF(t *testing.T) {
	testCases := []struct {
		desc     string
		r        func() *Bitset
		k        func() *Bitset
		expected []byte
	}{
		{
			desc: "it works#1",
			r: func() *Bitset {
				bits := []byte{
					0b11110000,
					0b10101010,
					0b11110000,
					0b10101010,
				}

				return BitsetFromBytes(bits)
			},
			k: func() *Bitset {
				bits := []byte{
					0b00011011,
					0b00000010,
					0b11101111,
					0b11111100,
					0b01110000,
					0b01110010,
				}

				return BitsetFromBytes(bits)
			},
			expected: []byte{
				0b00100011,
				0b01001010,
				0b10101001,
				0b10111011,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			enc := Encoder{}
			require.Equal(t, tC.expected, enc.f(tC.r(), tC.k()).Bits())
		})
	}
}
