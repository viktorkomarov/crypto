package des

const byteSize = 8

func or1(size int) int {
	if size/byteSize == 0 {
		return 1
	}

	return size / byteSize
}

type Bitset struct {
	buffer []byte
}

func BitsetFromSize(size int) *Bitset {
	return &Bitset{
		buffer: make([]byte, or1(size)),
	}
}

func BitsetFromBytes(bytes []byte) *Bitset {
	return &Bitset{
		buffer: bytes,
	}
}

func (b *Bitset) Nth(n int) byte {
	num := b.buffer[n/byteSize]
	return (num >> (byteSize - (n % byteSize) - 1)) & 1
}

// func (b *Bitset) SetVal(to int, val uint64) {
// 	appBlock := 0
// 	sz := to - b.cap() + 1
// 	if sz > 0 {
// 		appBlock = or1(uint64(sz))

// 	}
// 	for i := 0; i < appBlock*blockSize; i++ {
// 		b.buffer = append(b.buffer, 0)
// 	}

// 	if b.sz < to {
// 		b.sz = to
// 	}

// 	from := to / 64
// 	num := byteOrder.Uint64(b.buffer[from : from+blockSize])
// 	switch val {
// 	case 1:
// 		num |= (1 << (to % 64))
// 	case 0:
// 		num &= ^(1 << (to % 64))
// 	}
// 	byteOrder.PutUint64(b.buffer[from:from+blockSize], num)
// }

// func (b *Bitset) Size() int {
// 	return b.sz
// }

// func (b *Bitset) cap() int {
// 	return len(b.buffer) * 8
// }

// func (b *Bitset) Subset(from, to int) *Bitset {
// 	set := BitsetFromSize(uint64(to - from)) // check from > to
// 	for ; from < to; from++ {
// 		set.SetVal(from, b.Nth(from))
// 	}

// 	return set
// }

// func (b *Bitset) LeftRotate(shift int) *Bitset {
// 	shifted := make([]byte, 0)
// 	size := b.Size() - 1
// 	for i := 0; i < shift; i++ {
// 		shifted = append(shifted, byte(b.Nth(i)))
// 	}

// 	for i := 0; i < size-shift; i++ {
// 		b.SetVal(i, b.Nth(i+1))
// 	}

// 	for i, val := range shifted {
// 		b.SetVal(size-len(shifted)+i, uint64(val))
// 	}

// 	return b
// }

// func (b *Bitset) Append(bits *Bitset) {
// 	sz := b.Size()
// 	for i := 0; i < bits.Size(); i++ {
// 		b.SetVal(sz+i, bits.Nth(i))
// 	}
// }

// func (b *Bitset) Bits() []byte {
// 	bits := make([]byte, 0)
// 	for i := 0; i < b.Size()/8+1; i++ {
// 		bits = append(bits, b.buffer[i])
// 	}

// 	return bits
// }
