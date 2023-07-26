package slice

func Map[T comparable](src []T) map[T]struct{} {
	m := make(map[T]struct{}, len(src))
	for _, v := range src {
		m[v] = struct{}{}
	}
	return m
}
