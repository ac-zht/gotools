package mapping

func Keys[T comparable, A any](mp map[T]A) []T {
	keys := make([]T, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

func KeyExist[T comparable, A any](mp map[T]A, key T) bool {
	if _, ok := mp[key]; ok {
		return true
	}
	return false
}
