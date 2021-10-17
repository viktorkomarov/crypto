package aes

import (
	"errors"

	"github.com/viktorkomarov/crypto/bitset"
	"github.com/viktorkomarov/crypto/ffield"
)

type Cipher struct {
	words        []Word
	keyScheduler KeyScheduler
	nk           int
	nr           int
	nb           int
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
	c.keyScheduler = newKeyScheduuler128(key)
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

func (c Cipher) Encrypt(b []byte) []byte {
	set := bitset.SetFromBytes(b)

	for c.keyScheduler.Next() {
		subkey := make([]byte, 0)
		for _, w := range c.keyScheduler.Subkey() {
			tmp := bitset.SetFromUint32(uint32(w))
			subkey = append(subkey,
				tmp.Subset(0, 8).BuildUint8(),
				tmp.Subset(8, 16).BuildUint8(),
				tmp.Subset(16, 24).BuildUint8(),
				tmp.Subset(24, 32).BuildUint8(),
			)
		}

		set = set.XOR(bitset.SetFromBytes(subkey))
	}
}

func (c Cipher) subBytes(p *bitset.Set) *bitset.Set {
	for i := 0; i < 128; i++ {
		p.SetVal(i, sBox(p.Nth(i)))
	}

	return p
}

func (c Cipher) shiftRows(p *bitset.Set) *bitset.Set {
	bitSet1 := []byte{p.Nth(0), p.Nth(4), p.Nth(8), p.Nth(12)}
	bitSet2 := []byte{p.Nth(5), p.Nth(9), p.Nth(13), p.Nth(1)}
	bitSet3 := []byte{p.Nth(10), p.Nth(14), p.Nth(2), p.Nth(6)}
	bitSet4 := []byte{p.Nth(15), p.Nth(3), p.Nth(7), p.Nth(11)}

	bytes := make([]byte, 0, 15)
	for i := 0; i < 4; i++ {
		bytes = append(bytes, bitSet1[i], bitSet2[i], bitSet3[i], bitSet4[i])
	}

	return bitset.SetFromBytes(bytes)
}

func (c Cipher) mixColumns(p *bitset.Set)
