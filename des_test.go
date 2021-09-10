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
