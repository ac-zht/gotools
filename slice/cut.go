package slice

import "math"

// Cut 将切片分割为指定长度的多个子切片，如果元素未达到指定长度数量，则元素数量为准
func Cut[T any](src []T, length int) [][]T {
	cnt := 0
	ret := make([][]T, 0, int(math.Ceil(float64(len(src))/float64(length))))
	n := len(src)
	for cnt < n {
		left := cnt
		if cnt+length > n {
			length = n - cnt
		}
		cnt += length
		subset := make([]T, length)
		right := cnt + 1
		if right > n {
			right = n
		}
		copy(subset, src[left:right])
		ret = append(ret, subset)
	}
	return ret
}
