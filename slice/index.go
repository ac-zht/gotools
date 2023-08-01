package slice

// Index 返回值在切片中的第一个下标，未找到则返回-1
func Index[T comparable](val T, s []T) int {
	return IndexFunc[T](val, s, func(src, dst T) bool {
		return src == dst
	})
}

func IndexFunc[T comparable](val T, s []T, equal equalFunc[T]) int {
	for k, v := range s {
		if equal(val, v) {
			return k
		}
	}
	return -1
}

// LastIndex 返回值在切片中的最后一个下标，未找到则返回-1
func LastIndex[T comparable](val T, s []T) int {
	return LastIndexFunc[T](val, s, func(src, dst T) bool {
		return src == dst
	})
}

func LastIndexFunc[T comparable](val T, s []T, equal equalFunc[T]) int {
	last := len(s) - 1
	for i := last; i >= 0; i-- {
		if equal(val, s[i]) {
			return i
		}
	}
	return -1
}

// AllIndex 返回值在切片中的所有下标，未找到则返回空切片
func AllIndex[T comparable](val T, s []T) []int {
	return AllIndexFunc[T](val, s, func(src, dst T) bool {
		return src == dst
	})
}

func AllIndexFunc[T comparable](val T, s []T, equal equalFunc[T]) []int {
	res := make([]int, 0, len(s))
	for k, v := range s {
		if equal(val, v) {
			res = append(res, k)
		}
	}
	return res
}
