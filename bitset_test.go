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
		// {
		// 	desc:   "all 1",
		// 	buffer: []byte{255, 255, 255, 255},
		// 	expectedMap: map[int]uint64{
		// 		1:  1,
		// 		20: 1,
		// 		10: 1,
		// 		5:  1,
		// 		8:  1,
		// 		7:  1,
		// 	},
		// },
		// {
		// 	desc:   "all 0",
		// 	buffer: []byte{0, 0, 0, 0},
		// 	expectedMap: map[int]uint64{
		// 		1:  0,
		// 		20: 0,
		// 		10: 0,
		// 		5:  0,
		// 		8:  0,
		// 		7:  0,
		// 	},
		// },
		// {
		// 	desc:   "1 and 0",
		// 	buffer: []byte{5, 0, 10},
		// 	expectedMap: map[int]uint64{
		// 		0:  1,
		// 		1:  0,
		// 		2:  1,
		// 		4:  0,
		// 		8:  0,
		// 		11: 0,
		// 		17: 1,
		// 	},
		// },
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

// func TestSetVal(t *testing.T) {
// 	testCases := []struct {
// 		desc        string
// 		buffer      []byte
// 		setter      map[int]uint64
// 		expectedMap map[int]uint64
// 	}{
// 		{
// 			desc:   "#1",
// 			buffer: []byte{255},
// 			setter: map[int]uint64{
// 				0: 0, 1: 0, 2: 0,
// 				3: 0, 5: 0, 6: 0,
// 				7: 0,
// 			},
// 			expectedMap: map[int]uint64{
// 				0: 0, 1: 0, 2: 0,
// 				3: 0, 5: 0, 6: 0,
// 				7: 0, 4: 1,
// 			},
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			bits := BitsetFromBytes(tC.buffer)

// 			for key, val := range tC.setter {
// 				bits.SetVal(key, val)
// 			}

// 			for nth, val := range tC.expectedMap {
// 				require.Equal(t, val, bits.Nth(nth))
// 			}
// 		})
// 	}
// }

// func TestSetNewVal(t *testing.T) {
// 	testCases := []struct {
// 		desc     string
// 		size     uint64
// 		setter   []int
// 		expected int
// 	}{
// 		{
// 			desc:     "48 size",
// 			size:     48,
// 			setter:   []int{48, 59, 64, 72},
// 			expected: 72,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			bits := BitsetFromSize(tC.size)
// 			for _, val := range tC.setter {
// 				bits.SetVal(val, 1)
// 			}
// 			require.Equal(t, tC.expected, bits.Size())
// 		})
// 	}
// }

// func TestLeftRotate(t *testing.T) {
// 	testCases := []struct {
// 		desc     string
// 		args     []byte
// 		shift    int
// 		expected []byte
// 	}{
// 		{
// 			desc: "it works#1",
// 			args: []byte{
// 				0b11001101,
// 				0b00110011,
// 				0b11101110,
// 			},
// 			shift: 1,
// 			expected: []byte{
// 				0b01001101,
// 				0b00110011,
// 				0b11101111,
// 			},
// 		},
// 		// {
// 		// 	desc: "it works#2",
// 		// 	args: []byte{
// 		// 		0b11001101,
// 		// 		0b00110011,
// 		// 		0b11101110,
// 		// 	},
// 		// 	shift: 12,
// 		// 	expected: []byte{
// 		// 		0b00111110,
// 		// 		0b11101100,
// 		// 		0b011010011,
// 		// 	},
// 		// },
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			require.Equal(t, tC.expected, BitsetFromBytes(tC.args).LeftRotate(tC.shift).Bits())
// 		})
// 	}
// }
