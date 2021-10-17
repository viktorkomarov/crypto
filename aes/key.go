package aes

import (
	"github.com/viktorkomarov/crypto/bitset"
)

type Word uint32

type KeyScheduler interface {
	Next() bool
	Subkey() []Word
}

type keyScheduler128 struct {
	currRound int8
	words     []Word
}

func newKeyScheduuler128(key []byte) *keyScheduler128 {
	set := bitset.SetFromBytes(key)
	words := make([]Word, 4)
	for i := range words {
		wordSet := set.Subset(i*8, i*8+32)
		words = append(words, Word(wordSet.BuildUint32()))
	}
	return &keyScheduler128{
		words:     words,
		currRound: 0,
	}
}

func (k *keyScheduler128) Next() bool {
	return k.currRound <= 10
}

func (k *keyScheduler128) Subkey() []Word {
	k.currRound++
	curr := make([]Word, 4)
	copy(curr, k.words)
	k.words[0] = k.words[0] ^ k.g(k.words[3])
	k.words[1] = k.words[0] ^ k.words[1]
	k.words[2] = k.words[1] ^ k.words[2]
	k.words[3] = k.words[2] ^ k.words[3]
	return curr
}

func (k *keyScheduler128) g(w Word) Word {
	w = k.rotWord(w)
	w = k.subWord(w)
	return w
}

func (k *keyScheduler128) rotWord(w Word) Word {
	set := bitset.SetFromUint32(uint32(w))
	set = set.LeftRotate(8)
	return Word(set.BuildUint32())
}

func (k *keyScheduler128) subWord(w Word) Word {
	set := bitset.SetFromUint32(uint32(w))
	a0 := sBox(set.Subset(0, 8).BuildUint8())
	a1 := sBox(set.Subset(8, 16).BuildUint8())
	a2 := sBox(set.Subset(16, 24).BuildUint8())
	a3 := sBox(set.Subset(24, 32).BuildUint8())
	return Word(bitset.SetFromBytes([]byte{a0, a1, a2, a3}).BuildUint32())
}

func (k *keyScheduler128) rCon(w Word) Word {
	c := []uint32{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x1B, 0x36}
	return Word(uint32(w) ^ c[k.currRound-1])
}
