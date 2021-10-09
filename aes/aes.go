package aes

import (
	"errors"

	"github.com/viktorkomarov/crypto/bitset"
	"github.com/viktorkomarov/crypto/ffield"
)

type Cipher struct {
	words []Word
	nk    int
	nr    int
	nb    int
}

var ErrInvalidKeySize = errors.New("invalid key size")

func keyRound(sz int) (int, error) {
	switch sz {
	case 128:
		return 10, nil
	case 192:
		return 12, nil
	case 256:
		return 14, nil
	default:
		return 0, ErrInvalidKeySize
	}
}

func NewCipher(key []byte) (c Cipher, err error) {
	sz := len(key) * 8
	c.key = bitset.SetFromBytes(key)
	c.nk = sz / 32
	c.nr, err = keyRound(sz)
	if err != nil {
		return c, err
	}
	c.nb = sz / 32
	return
}

func sBox(a uint8) uint8 {
	invr := uint8(ffield.InvrGF8(uint64(a)))
	b := bitset.SetFromUint8(invr)
	c := bitset.SetFromUint8(0b01100011)
	result := bitset.SetFromUint8(0)
	for i := 0; i < 8; i++ {
		result.SetVal(i, b.Nth(i)^b.Nth((i+4)%8)^b.Nth((i+5)%8)^b.Nth((i+6)%8)^b.Nth((i+7)%8)^c.Nth(i))
	}
	return result.BuildUint8()
}
