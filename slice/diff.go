package slice

// Diff 切片src与dst的差集
func Diff[T comparable](src, dst []T) []T {
	srcMap := Map[T](src)
	for _, v := range dst {
		delete(srcMap, v)
	}
	diff := make([]T, 0, len(srcMap))
	for k, _ := range srcMap {
		diff = append(diff, k)
	}
	return diff
}

// SymmetricDiff 切片src与dst的对称差集
func SymmetricDiff[T comparable](src, dst []T) []T {
	srcMap := Map[T](src)
	dstMap := Map[T](dst)
	for _, v := range dst {
		if _, ok := srcMap[v]; ok {
			delete(srcMap, v)
			delete(dstMap, v)
		}
	}
	for k, _ := range dstMap {
		srcMap[k] = struct{}{}
	}
	sd := make([]T, 0, len(srcMap))
	for k, _ := range srcMap {
		sd = append(sd, k)
	}
	return sd
}
