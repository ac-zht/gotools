package slice

import "gotools"

// Max 切片中最大值
func Max[T gotools.RealNumber](src []T) T {
	res := src[0]
	for _, v := range src {
		if v > res {
			res = v
		}
	}
	return res
}

// Min 切片中最小值
func Min[T gotools.RealNumber](src []T) T {
	res := src[0]
	for _, v := range src {
		if v < res {
			res = v
		}
	}
	return res
}

// Sum 切片中值的总和
func Sum[T gotools.Number](src []T) T {
	var res T
	for _, v := range src {
		res += v
	}
	return res
}
