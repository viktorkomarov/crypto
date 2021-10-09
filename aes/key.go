package aes

import (
	"github.com/viktorkomarov/crypto/bitset"
	"github.com/viktorkomarov/crypto/ffield"
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
	return k.currRound < 10
}

func (k *keyScheduler128) Subkey() []Word {
	curr := make([]Word, 4)
	copy(curr, k.words)
	k.words[0] = k.words[0] ^ k.g(k.words[3])
	k.words[1] = k.words[0] ^ k.words[1]
	k.words[2] = k.words[1] ^ k.words[2]
	k.words[3] = k.words[2] ^ k.words[3]
	return curr
}

func (k *keyScheduler128) g(w Word) Word {
	w := k.rotWord(w)
}

func (k *keyScheduler128) rotWord(w Word) Word {
	set := bitset.SetFromUin32(uint32(w))
	set = set.LeftRotate(8)
	return Word(set.BuildUint32())
}

func (k *keyScheduler128) sBox(w Word) Word {
	invr := ffield.InvrGF8(uint64(w))
}
