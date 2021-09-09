package des

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitsetNth(t *testing.T) {
	testCases := []struct {
		desc        string
		buffer      []byte
		expectedMap map[int]byte
	}{
		{
			desc:   "all 1",
			buffer: []byte{255, 255, 255, 255},
			expectedMap: map[int]byte{
				1: 1, 20: 1,
				10: 1, 5: 1,
				8: 1, 7: 1,
			},
		},
		{
			desc:   "all 0",
			buffer: []byte{0, 0, 0, 0},
			expectedMap: map[int]byte{
				1: 0, 20: 0,
				10: 0, 5: 0,
				8: 0, 7: 0,
			},
		},
		{
			desc: "1 and 0",
			buffer: []byte{
				0b00000101,
				0b00000000,
				0b00001010,
			},
			expectedMap: map[int]byte{
				0: 0, 1: 0, 5: 1,
				7: 1, 8: 0, 11: 0,
				20: 1, 21: 0, 22: 1,
			},
		},
		{
			desc: "step by step",
			buffer: []byte{
				0b11001101,
				0b00110011,
				0b11101110,
			},
			expectedMap: map[int]byte{
				0: 1, 1: 1, 2: 0, 3: 0, 4: 1, 5: 1, 6: 0, 7: 1,
				8: 0, 9: 0, 10: 1, 11: 1, 12: 0, 13: 0, 14: 1, 15: 1,
				16: 1, 17: 1, 18: 1, 19: 0, 20: 1, 21: 1, 22: 1, 23: 0,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := BitsetFromBytes(tC.buffer)

			for nth, val := range tC.expectedMap {
				require.Equal(t, val, bits.Nth(nth))
			}
		})
	}
}

func TestSetVal(t *testing.T) {
	testCases := []struct {
		desc        string
		buffer      []byte
		setter      map[int]byte
		expectedMap map[int]byte
	}{
		{
			desc:   "#1",
			buffer: []byte{255},
			setter: map[int]byte{
				0: 0, 1: 0, 2: 0,
				3: 0, 5: 0, 6: 0,
				7: 0,
			},
			expectedMap: map[int]byte{
				0: 0, 1: 0, 2: 0,
				3: 0, 5: 0, 6: 0,
				7: 0, 4: 1,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := BitsetFromBytes(tC.buffer)

			for key, val := range tC.setter {
				bits.SetVal(key, val)
			}

			for nth, val := range tC.expectedMap {
				require.Equal(t, val, bits.Nth(nth))
			}
		})
	}
}

func TestSetNewVal(t *testing.T) {
	testCases := []struct {
		desc     string
		size     int
		setter   []int
		expected int
	}{
		{
			desc:     "48 size",
			size:     48,
			setter:   []int{48, 59, 64, 72},
			expected: 72,
		},
		{
			desc:     "10 size",
			size:     10,
			setter:   []int{11, 12, 13, 14, 15, 16, 17, 18},
			expected: 18,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := BitsetFromSize(tC.size)
			for _, val := range tC.setter {
				bits.SetVal(val, 1)
			}
			require.Equal(t, tC.expected, bits.Size())
		})
	}
}

func TestLeftRotate(t *testing.T) {
	testCases := []struct {
		desc     string
		args     []byte
		shift    int
		expected []byte
	}{
		{
			desc: "it works#1",
			args: []byte{
				0b11001101,
				0b00110011,
				0b11101110,
			},
			shift: 1,
			expected: []byte{
				0b10011010,
				0b01100111,
				0b11011101,
			},
		},
		{
			desc: "it works#2",
			args: []byte{
				0b11001101,
				0b00110011,
				0b11101110,
			},
			shift: 4,
			expected: []byte{
				0b11010011,
				0b00111110,
				0b11101100,
			},
		},
		{
			desc: "it works#3",
			args: []byte{
				0b11001000,
			},
			shift: 2,
			expected: []byte{
				0b00100011,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bits := BitsetFromBytes(tC.args)
			require.Equal(t, tC.expected, bits.LeftRotate(tC.shift).Bits())
		})
	}
}
