package des

import (
	"crypto/des"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/viktorkomarov/crypto/bitset"
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
			require.Equal(t, tC.expected, NewCipher(tC.key).key.Bits())
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
			bits := NewCipher(key)
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
		r        func() *bitset.Set
		k        func() *bitset.Set
		expected []byte
	}{
		{
			desc: "it works#1",
			r: func() *bitset.Set {
				bits := []byte{
					0b11110000,
					0b10101010,
					0b11110000,
					0b10101010,
				}

				return bitset.SetFromBytes(bits)
			},
			k: func() *bitset.Set {
				bits := []byte{
					0b00011011,
					0b00000010,
					0b11101111,
					0b11111100,
					0b01110000,
					0b01110010,
				}

				return bitset.SetFromBytes(bits)
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
			enc := Cipher{}
			require.Equal(t, tC.expected, enc.f(tC.r(), tC.k()).Bits())
		})
	}
}

func TestEncoder(t *testing.T) {
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
	msg := []byte{
		0b00000001,
		0b00100011,
		0b01000101,
		0b01100111,
		0b10001001,
		0b10101011,
		0b11001101,
		0b11101111,
	}

	destTest, err := des.NewCipher(key)
	require.NoError(t, err)
	desNature := make([]byte, 8)
	destTest.Encrypt(desNature, msg)
	encoder := NewCipher(key)
	require.Equal(t, desNature, encoder.Encrypt(msg))
}

/*
goos: linux
goarch: amd64
pkg: github.com/viktorkomarov/des
BenchmarkStdDESEncoder-16    	  550294	      1838 ns/op	     264 B/op	       4 allocs/op
*/
func BenchmarkStdDESEncoder(b *testing.B) {
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
	msg := []byte{
		0b00000001,
		0b00100011,
		0b01000101,
		0b01100111,
		0b10001001,
		0b10101011,
		0b11001101,
		0b11101111,
	}

	for i := 0; i < b.N; i++ {
		destTest, err := des.NewCipher(key)
		if err != nil {
			b.Fail()
		}

		desNature := make([]byte, 8)
		destTest.Encrypt(desNature, msg)
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/viktorkomarov/des
BenchmarkMyDESEncoder-16    	   10821	    111551 ns/op	   10320 B/op	     615 allocs/op
*/
func BenchmarkMyDESEncoder(b *testing.B) {
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
	msg := []byte{
		0b00000001,
		0b00100011,
		0b01000101,
		0b01100111,
		0b10001001,
		0b10101011,
		0b11001101,
		0b11101111,
	}

	for i := 0; i < b.N; i++ {
		NewCipher(key).Encrypt(msg)
	}
}

func TestDecrypt(t *testing.T) {
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
	msg := []byte{
		0b00000001,
		0b00100011,
		0b01000101,
		0b01100111,
		0b10001001,
		0b10101011,
		0b11001101,
		0b11101111,
	}

	d, err := des.NewCipher(key)
	require.NoError(t, err)
	enc := make([]byte, 8)
	d.Encrypt(enc, msg)
	require.Equal(t, msg, NewCipher(key).Decrypt(enc))
}
