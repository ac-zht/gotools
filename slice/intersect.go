package slice

// Intersect 两切片的交集
func Intersect[T comparable](src, dst []T) []T {
	m := Map[T](src)
	l := len(src)
	if len(dst) < l {
		l = len(dst)
	}
	intersect := make([]T, 0, l)
	for _, v := range dst {
		if _, ok := m[v]; ok {
			intersect = append(intersect, v)
		}
	}
	return intersect
}
