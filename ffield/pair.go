package ffield

type Pair struct {
	L, R uint64
}

func NewPair(l, r uint64) Pair {
	if r < l {
		l, r = r, l
	}

	return Pair{
		L: l,
		R: r,
	}
}
