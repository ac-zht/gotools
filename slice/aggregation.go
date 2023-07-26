package slice

import "gotools"

func Max[T gotools.RealNumber](src []T) T {
	res := src[0]
	for _, v := range src {
		if v > res {
			res = v
		}
	}
	return res
}

func Min[T gotools.RealNumber](src []T) T {
	res := src[0]
	for _, v := range src {
		if v < res {
			res = v
		}
	}
	return res
}

func Sum[T gotools.Number](src []T) T {
	var res T
	for _, v := range src {
		res += v
	}
	return res
}
