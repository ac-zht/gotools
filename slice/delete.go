package slice

func Delete[T any](index int, s []T) []T {
	if index < 0 || index >= len(s) {
		return s
	}
	return FilterDelete(s, func(key int, value T) bool {
		if key == index {
			return true
		}
		return false
	})
}

func FilterDelete[T any](s []T, filter func(key int, value T) bool) []T {
	i := 0
	for k, v := range s {
		if filter(k, v) {
			continue
		}
		s[i] = v
		i++
	}
	return s[:i]
}
