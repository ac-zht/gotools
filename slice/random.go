package slice

import "gotools/random"

func RandomIndexes[T any]() {

}

// WeightRandomIndex 按整型权重随机返回切片节点
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
