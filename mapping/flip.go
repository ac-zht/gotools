package mapping

// Flip 交换map中的键和值
// 由于map无序所以无法确认交换后键的最终覆盖值，只用于没有重复值的map
func Flip[T1 comparable, T2 comparable](src map[T1]T2) map[T2]T1 {
	ret := make(map[T2]T1, len(src))
	for k, v := range src {
		ret[v] = k
	}
	return ret
}
