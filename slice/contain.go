package slice

// Contain 切片中是否包含某元素
func Contain[T comparable](s []T, val T) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}
