package mapping

func Keys[T comparable, A any](mp map[T]A) []T {
	keys := make([]T, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}
