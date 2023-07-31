package slice

func Intersection[T comparable](src, dst []T) []T {
	m := Map[T](src)
	l := len(src)
	if len(dst) < l {
		l = len(dst)
	}
	intersection := make([]T, 0, l)
	for _, v := range dst {
		if _, ok := m[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}
