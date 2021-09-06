package des

import "encoding/binary"

var byteOrder = binary.BigEndian

const blockSize = 8

type Bitset struct {
	sz     int
	buffer []byte
}

func or1(size uint64) int {
	if size/64 == 0 {
		return blockSize
	}
	return int(size / 64)
}

func BitsetFromSize(size uint64) *Bitset {
	return &Bitset{
		buffer: make([]byte, or1(size)),
	}
}

func BitsetFromBytes(bytes []byte) *Bitset {
	bits := BitsetFromSize(uint64(len(bytes) * 8))
	for i := range bytes {
		for j := 0; j < 8; j++ {
			num := (bytes[i] >> j) & 1
			bits.SetVal(i*8+j, uint64(num))
		}
	}

	return bits
}

func (b *Bitset) Nth(n int) uint64 {
	from := n / 64
	return (byteOrder.Uint64(b.buffer[from:from+blockSize]) >> (n % 64)) & 1
}

func (b *Bitset) SetVal(to int, val uint64) {
	appBlock := 0
	sz := to - b.Size() + 1
	if sz > 0 {
		appBlock = or1(uint64(sz))

	}
	for i := 0; i < appBlock*blockSize; i++ {
		b.buffer = append(b.buffer, 0)
	}

	if b.sz < to {
		b.sz = to
	}

	from := to / 64
	num := byteOrder.Uint64(b.buffer[from : from+blockSize])
	switch val {
	case 1:
		num |= (1 << (to % 64))
	case 0:
		num &= ^(1 << (to % 64))
	}
	byteOrder.PutUint64(b.buffer[from:from+blockSize], num)
}

func (b *Bitset) Size() int {
	return b.sz
}

func (b *Bitset) Subset(from, to int) *Bitset {
	set := BitsetFromSize(uint64(to - from)) // check from > to
	for ; from < to; from++ {
		set.SetVal(from, b.Nth(from))
	}

	return set
}

func (b *Bitset) LeftRotate(shift int) *Bitset {
	shifted := make([]byte, 0)
	size := b.Size() - 1
	for i := 0; i < shift; i++ {
		shifted = append(shifted, byte(b.Nth(size-i)))
	}

	for i := 0; i < size-shift; i++ {
		b.SetVal(size-i, b.Nth(size-i))
	}

	for i, val := range shifted {
		b.SetVal(i, uint64(val))
	}

	return b
}

func (b *Bitset) Append(bits *Bitset) {
	sz := b.Size()
	for i := 0; i < bits.Size(); i++ {
		b.SetVal(sz+i, bits.Nth(i))
	}
}

func (b *Bitset) Bits() []byte {
	bits := make([]byte, len(b.buffer))
	copy(bits, b.buffer)
	return bits
}
