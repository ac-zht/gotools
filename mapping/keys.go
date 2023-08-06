package mapping

// Keys 获取Map所有键名
func Keys[T comparable, A any](mp map[T]A) []T {
	keys := make([]T, 0, len(mp))
	for k := range mp {
		keys = append(keys, k)
	}
	return keys
}

// KeyExist 判断Map中键名是否存在
func KeyExist[T comparable, A any](mp map[T]A, key T) bool {
	if _, ok := mp[key]; ok {
		return true
	}
	return false
}

// KeysWithValue 获取指定值对应的所有键名
func KeysWithValue[T1 comparable, T2 comparable](src map[T1]T2, val T2) []T1 {
	keys := make([]T1, 0, len(src))
	for k, v := range src {
		if v == val {
			keys = append(keys, k)
		}
	}
	return keys
}
