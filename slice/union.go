package slice

// Union 两切片的全集
func Union[T comparable](s1 []T, s2 []T) []T {
	m1 := Map[T](s1)
	s := make([]T, 0, len(s1)+len(s2))
	s = s1
	for _, v := range s2 {
		if _, ok := m1[v]; ok {
			continue
		}
		s = append(s, v)
	}
	return s
}
