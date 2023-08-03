package slice

// Unique 切片去重
func Unique[T comparable](src []T) []T {
	srcMap := Map[T](src)
	uniq := make([]T, 0, len(srcMap))
	for k, _ := range srcMap {
		uniq = append(uniq, k)
	}
	return uniq
}
