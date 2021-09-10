package des

import (
	"fmt"
	"strings"
)

const byteSize = 8

func or1(size int) int {
	if size/byteSize <= 0 {
		return 1
	}

	return size/byteSize + 1
}

// think about LittleEndian
type Bitset struct {
	sz     int
	buffer []byte
}

func BitsetFromSize(size int) *Bitset {
	return &Bitset{
		sz:     size,
		buffer: make([]byte, or1(size)),
	}
}

func BitsetFromBytes(bytes []byte) *Bitset {
	return &Bitset{
		sz:     len(bytes) * byteSize,
		buffer: bytes,
	}
}

func (b *Bitset) Nth(n int) byte {
	num := b.buffer[n/byteSize]
	return (num >> (byteSize - (n % byteSize) - 1)) & 1
}

func (b *Bitset) nthOr0(n int) byte {
	if n >= b.Size() {
		return 0
	}

	return b.Nth(n)
}

func (b *Bitset) SetVal(to int, val byte) {
	if b.sz < to {
		b.sz = to
	}

	appBlock := 0
	sz := to - b.cap()
	if sz >= 0 {
		appBlock = or1(sz)

	}

	for i := 0; i < appBlock; i++ {
		b.buffer = append(b.buffer, 0)
	}

	switch val {
	case 1:
		b.buffer[to/byteSize] |= (1 << (byteSize - 1 - (to % byteSize)))
	case 0:
		b.buffer[to/byteSize] &= ^(1 << (byteSize - 1 - (to % byteSize)))
	}
}

func (b *Bitset) Size() int {
	return b.sz
}

func (b *Bitset) cap() int {
	return len(b.buffer) * 8
}

func (b *Bitset) Subset(from, to int) *Bitset {
	set := BitsetFromSize(to - from) // check from > to
	i := 0
	for ; from < to; from++ {
		set.SetVal(i, b.Nth(from))
		i += 1
	}

	return set
}

func (b *Bitset) LeftRotate(shift int) *Bitset {
	shift %= b.sz
	bits := BitsetFromSize(b.sz)

	shifted := make([]byte, 0)
	for i := 0; i < shift; i++ {
		shifted = append(shifted, b.Nth(i))
	}

	for i := 0; i < b.sz-shift; i++ {
		bits.SetVal(i, b.Nth(i+shift))
	}

	for i := 0; i < shift; i++ {
		bits.SetVal(b.sz-shift+i, shifted[i])
	}

	return bits
}

func (b *Bitset) Bits() []byte {
	bits := make([]byte, b.sz/8)
	copy(bits, b.buffer)
	return bits
}

func (b *Bitset) String() string {
	var builder strings.Builder

	for i := 0; i < b.sz; i++ {
		builder.WriteString(fmt.Sprintf("%d", b.Nth(i)))
	}

	return builder.String()
}

func (b *Bitset) XOR(bits *Bitset) *Bitset {
	max, min := b, bits
	if max.Size() < min.Size() {
		max, min = min, max
	}

	result := BitsetFromSize(max.Size())
	for i := 0; i < max.Size(); i++ {
		result.SetVal(i, (max.Nth(i) ^ min.nthOr0(i)))
	}

	return result
}
