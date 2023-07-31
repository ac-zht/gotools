package slice

func Diff[T comparable](src, dst []T) []T {
	m := Map[T](src)
	l := len(src)
	if len(dst) < l {
		l = len(dst)
	}
	diff := make([]T, 0, l)
	for _, v := range dst {
		if _, ok := m[v]; ok {
			diff = append(diff, v)
		}
	}
	return diff
}
