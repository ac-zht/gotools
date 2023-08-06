package slice

// Combine 创建一个Map，一个数组的值作为Map键名，另一数组的值作为Map值
func Combine[T comparable, A any](keys []T, values []A) map[T]A {
	mp := make(map[T]A, len(keys))
	for k, v := range keys {
		mp[v] = values[k]
	}
	return mp
}
