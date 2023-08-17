package slice

import (
	"github.com/zht-account/gotools"
	"github.com/zht-account/gotools/random"
)

// RandomIndexes 随机得到切片中指定数量的下标集
func RandomIndexes[T any](src []T, n int) ([]int, error) {
	if len(src) <= 0 {
		return nil, gotools.ErrSliceIsEmpty
	}
	if n > len(src) {
		return nil, gotools.ErrSliceLengthNotEnough
	}
	ret := make([]int, 0, n)
	cnt := 0
	for cnt <= n {
		m := random.RandInt(0, len(src)-1)
		if Contain[int](ret, m) {
			continue
		}
		ret = append(ret, m)
		cnt++
	}
	return ret, nil
}

// WeightRandomIndex 按整型权重随机返回切片某个下标
func WeightRandomIndex(src []int) int {
	sum := Sum[int](src)
	n := random.RandInt(1, sum)
	var limit int
	for k, v := range src {
		limit += v
		if n <= limit {
			return k
		}
	}
	return 0
}
