package aes

import (
	"errors"

	"github.com/viktorkomarov/crypto/bitset"
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

func sBox(b uint8) uint8 {

}
