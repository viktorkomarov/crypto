package des

type Bitset struct {
	buffer []byte
}

func or1(size int) int {
	if size/8 == 0 {
		return 1
	}
	return size / 8
}

func BitsetFromSize(size int) *Bitset {
	return &Bitset{
		buffer: make([]byte, or1(size)),
	}
}

func BitsetFromBytes(bytes []byte) *Bitset {
	buffer := make([]byte, len(bytes))
	copy(buffer, bytes)
	return &Bitset{
		buffer: buffer,
	}
}

func (b *Bitset) Set(bit int) {
	b.SetVal(bit, 1)
}

func (b *Bitset) Unset(bit int) {
	b.SetVal(bit, 0)
}

func (b *Bitset) Nth(bit int) byte {
	return (b.buffer[bit/8] >> (bit % 8)) & 1
}

func (b *Bitset) SetVal(to int, val byte) {
	appBlock := 0
	sz := to - b.Size() + 1
	if sz > 0 {
		appBlock = or1(sz)

	}
	for i := 0; i < appBlock; i++ {
		b.buffer = append(b.buffer, 0)
	}

	switch val {
	case 1:
		b.buffer[to/8] |= (1 << (to % 8))
	case 0:
		b.buffer[to/8] &= ^(1 << (to % 8))
	}
}

func (b *Bitset) Size() int {
	return 8 * len(b.buffer)
}

func (b *Bitset) Subset(from, to int) *Bitset {
	bits := BitsetFromSize(to - from)

	for ; from < to; from++ {
		bits.SetVal(from, b.Nth(from))
	}

	return bits
}

func (b *Bitset) LeftRotate(shift int) *Bitset {
	shifted := make([]byte, 0)
	size := b.Size() - 1
	for i := 0; i < shift; i++ {
		shifted = append(shifted, b.Nth(size-i))
	}

	for i := 0; i < size-shift; i++ {
		b.SetVal(size-i, b.Nth(size-i))
	}

	for i, val := range shifted {
		b.SetVal(i, val)
	}

	return b
}

func (b *Bitset) Append(bits *Bitset) {
	for i := 0; i < bits.Size(); i++ {
		b.SetVal(b.Size()+i, bits.Nth(i))
	}
}

func (b *Bitset) Bits() []byte {
	bits := make([]byte, len(b.buffer))
	copy(bits, b.buffer)
	return bits
}
