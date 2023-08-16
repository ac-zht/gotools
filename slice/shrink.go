package slice

func CalcCapacity(l, c int) (int, bool) {
	if c > 64 && c <= 2048 && (c/l >= 4) {
		return c / 2, true
	} else if c > 2048 && (c/l >= 2) {
		return int(float64(c) * 0.625), true
	}
	return c, false
}

func Shrink[T any](src []T) []T {
	l, c := len(src), cap(src)
	capacity, changed := CalcCapacity(l, c)
	if !changed {
		return src
	}
	dst := make([]T, len(src), capacity)
	copy(dst, src)
	return dst
}
