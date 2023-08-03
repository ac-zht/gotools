package slice

// CountValues 统计切片中每个元素的数量
func CountValues[T comparable](src []T) map[T]int {
	valCntMap := make(map[T]int, len(src))
	for _, v := range src {
		valCntMap[v]++
	}
	return valCntMap
}
