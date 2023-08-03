package slice

// Fill 用指定值填充切片
func Fill[T any](start, cnt int, val T) []T {
	ret := make([]T, start+cnt)
	for k, _ := range ret {
		if k < start {
			continue
		}
		ret[k] = val
	}
	return ret
}
