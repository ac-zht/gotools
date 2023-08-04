package slice

// Merge 多个切片合并
func Merge[T any](many ...[]T) []T {
	ret := make([]T, 0, len(many[0]))
	for _, v := range many {
		ret = append(ret, v...)
	}
	return ret
}
