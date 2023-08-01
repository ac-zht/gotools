package slice

// Reverse 用新切片接收翻转值
func Reverse[T any](src []T) []T {
	ret := make([]T, len(src))
	j := 0
	for i := len(src) - 1; i >= 0; i-- {
		ret[j] = src[i]
		j++
	}
	return ret
}

// ReverseSelf 在自身切片上进行翻转
func ReverseSelf[T any](src []T) {
	for i, j := 0, len(src)-1; i < j; i, j = i+1, j-1 {
		src[i], src[j] = src[j], src[i]
	}
}
