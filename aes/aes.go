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
	if c.keyScheduler.Next() {
		set = c.addKeyLayer(set)
	}

	for c.keyScheduler.Next() {
		set = c.addKeyLayer(c.mixColumns(c.shiftRows(c.subBytes(set))))
	}

	return set.Bits()
}

func (c Cipher) subBytes(p *bitset.Set) *bitset.Set {
	for i := 0; i < 128; i++ {
		p.SetVal(i, sBox(p.Nth(i)))
	}

	return p
}

func buildSetFromMatrix(b ...[]byte) *bitset.Set {
	bytes := make([]byte, 0, 15)
	for i := 0; i < 4; i++ {
		bytes = append(bytes, b[0][i], b[1][i], b[2][i], b[3][i])
	}

	return bitset.SetFromBytes(bytes)
}

func (c Cipher) shiftRows(p *bitset.Set) *bitset.Set {
	bitSet1 := []byte{p.Nth(0), p.Nth(4), p.Nth(8), p.Nth(12)}
	bitSet2 := []byte{p.Nth(5), p.Nth(9), p.Nth(13), p.Nth(1)}
	bitSet3 := []byte{p.Nth(10), p.Nth(14), p.Nth(2), p.Nth(6)}
	bitSet4 := []byte{p.Nth(15), p.Nth(3), p.Nth(7), p.Nth(11)}

	return buildSetFromMatrix(bitSet1, bitSet2, bitSet3, bitSet4)
}

func (c Cipher) mixColumns(p *bitset.Set) *bitset.Set {
	col1 := []byte{
		byte(ffield.MulGF8(uint64(p.Nth(0)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(1)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(2)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(3)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(0)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(1)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(2)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(3)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(0)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(1)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(2)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(3)), uint64(3))),
		byte(ffield.MulGF8(uint64(p.Nth(0)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(1)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(2)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(1)), uint64(2))),
	}
	col2 := []byte{
		byte(ffield.MulGF8(uint64(p.Nth(4)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(5)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(6)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(7)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(4)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(5)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(6)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(7)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(4)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(5)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(6)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(7)), uint64(3))),
		byte(ffield.MulGF8(uint64(p.Nth(4)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(5)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(6)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(7)), uint64(2))),
	}
	col3 := []byte{
		byte(ffield.MulGF8(uint64(p.Nth(8)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(9)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(10)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(11)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(8)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(9)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(10)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(11)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(8)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(9)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(10)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(11)), uint64(3))),
		byte(ffield.MulGF8(uint64(p.Nth(8)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(9)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(10)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(11)), uint64(2))),
	}
	col4 := []byte{
		byte(ffield.MulGF8(uint64(p.Nth(12)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(13)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(14)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(15)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(12)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(13)), uint64(2)) ^ ffield.MulGF8(uint64(p.Nth(14)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(15)), uint64(1))),
		byte(ffield.MulGF8(uint64(p.Nth(12)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(13)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(14)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(15)), uint64(3))),
		byte(ffield.MulGF8(uint64(p.Nth(12)), uint64(3)) ^ ffield.MulGF8(uint64(p.Nth(13)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(14)), uint64(1)) ^ ffield.MulGF8(uint64(p.Nth(15)), uint64(2))),
	}

	return buildSetFromMatrix(col1, col2, col3, col4)
}

func (c Cipher) addKeyLayer(p *bitset.Set) *bitset.Set {
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

	return p.XOR(bitset.SetFromBytes(subkey))
}
