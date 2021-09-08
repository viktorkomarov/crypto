package des

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
		buffer: make([]byte, or1(size-1)),
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
	for ; from < to; from++ {
		set.SetVal(from, b.Nth(from))
	}

	return set
}

func (b *Bitset) LeftRotate(shift int) {
	shifted := make([]byte, 0)
	for i := 0; i < shift; i++ {
		shifted = append(shifted, b.Nth(i))
	}

	for i := 0; i < b.sz-shift; i++ {
		b.SetVal(i, b.Nth(i+1))
	}

	for i := 0; i < shift; i++ {
		b.SetVal(b.sz-shift+i, shifted[i])
	}
}

func (b *Bitset) Append(bits *Bitset) {
	sz := b.Size()
	for i := 0; i < bits.Size(); i++ {
		b.SetVal(sz+i, bits.Nth(i))
	}
}

func (b *Bitset) Bits() []byte {
	return b.buffer
}
