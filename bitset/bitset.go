package bitset

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
type Set struct {
	sz     int
	buffer []byte
}

func SetFromSize(size int) *Set {
	return &Set{
		sz:     size,
		buffer: make([]byte, or1(size)),
	}
}

func SetFromBytes(bytes []byte) *Set {
	return &Set{
		sz:     len(bytes) * byteSize,
		buffer: bytes,
	}
}

func (b *Set) Nth(n int) byte {
	num := b.buffer[n/byteSize]
	return (num >> (byteSize - (n % byteSize) - 1)) & 1
}

func (b *Set) nthOr0(n int) byte {
	if n >= b.Size() {
		return 0
	}

	return b.Nth(n)
}

func (b *Set) SetVal(to int, val byte) {
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

func (b *Set) Size() int {
	return b.sz
}

func (b *Set) cap() int {
	return len(b.buffer) * 8
}

func (b *Set) Subset(from, to int) *Set {
	set := SetFromSize(to - from) // check from > to
	i := 0
	for ; from < to; from++ {
		set.SetVal(i, b.Nth(from))
		i += 1
	}

	return set
}

func (b *Set) LeftRotate(shift int) *Set {
	shift %= b.sz
	bits := SetFromSize(b.sz)

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

func (b *Set) Bits() []byte {
	bits := make([]byte, b.sz/8)
	copy(bits, b.buffer)
	return bits
}

func (b *Set) String() string {
	var builder strings.Builder

	for i := 0; i < b.sz; i++ {
		builder.WriteString(fmt.Sprintf("%d", b.Nth(i)))
	}

	return builder.String()
}

func (b *Set) XOR(bits *Set) *Set {
	max, min := b, bits
	if max.Size() < min.Size() {
		max, min = min, max
	}

	result := SetFromSize(max.Size())
	for i := 0; i < max.Size(); i++ {
		result.SetVal(i, (max.Nth(i) ^ min.nthOr0(i)))
	}

	return result
}
