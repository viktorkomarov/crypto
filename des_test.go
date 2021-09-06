package des

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEncoder(t *testing.T) {
	key := []byte{
		0b11001000,
		0b00101100,
		0b11101010,
		0b10011110,
		0b11011001,
		0b00111101,
		0b11111011,
		0b10001111,
	}

	encoder := NewEncoder(key)
	expected := []byte{
		0b00001111,
		0b00110011,
		0b01010101,
		0b10101111,
		0b01101010,
		0b11100110,
		0b11110001,
	}
	require.Equal(t, expected, encoder.key.Bits())
}

func TestKey(t *testing.T) {
	key := []byte{
		0b11001000,
		0b00101100,
		0b11101010,
		0b10011110,
		0b11011001,
		0b00111101,
		0b11111011,
		0b10001111,
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
				0b11011000,
				0b01000000,
				0b11110111,
				0b00111111,
				0b00001110,
				0b01001110,
			},
		},
		// {
		// 	desc:  "round 2",
		// 	round: 1,
		// 	expected: []byte{
		// 		0b01111001,
		// 		0b10101110,
		// 		0b11011001,
		// 		0b11011011,
		// 		0b11001001,
		// 		0b11100101,
		// 	},
		// },
		// {
		// 	desc:  "last round",
		// 	round: 15,
		// 	expected: []byte{
		// 		0b11001011,
		// 		0b00111101,
		// 		0b10001011,
		// 		0b00001110,
		// 		0b00010111,
		// 		0b11110101,
		// 	},
		// },
		// {
		// 	desc:  "last round",
		// 	round: 15,
		// 	expected: []byte{
		// 		0b11001011,
		// 		0b00111101,
		// 		0b10001011,
		// 		0b00001110,
		// 		0b00010111,
		// 		0b11110101,
		// 	},
		// },
		// {
		// 	desc:  "round 8",
		// 	round: 7,
		// 	expected: []byte{
		// 		0b11110111,
		// 		0b10001010,
		// 		0b00111010,
		// 		0b11000001,
		// 		0b00111011,
		// 		0b11111011,
		// 	},
		// },
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := NewEncoder(key)
			require.Equal(t, tC.expected, bits.round(tC.round))
		})
	}
}
