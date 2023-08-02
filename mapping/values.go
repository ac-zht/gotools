package mapping

func Values[T comparable, A any](mp map[T]A) []A {
	values := make([]A, 0, len(mp))
	for _, v := range mp {
		values = append(values, v)
	}
	return values
}
