package slice

// Column 返回map切片中的指定列值
func Column[T comparable, A any](src []map[T]A, colKey T) []A {
	ret := make([]A, len(src))
	for key, val := range src {
		if v, ok := val[colKey]; ok {
			ret[key] = v
		}
	}
	return ret
}

// ColumnWithFilterNotExist 返回map切片中的指定列值, 过滤掉不存在的列
func ColumnWithFilterNotExist[T comparable, A any](src []map[T]A, colKey T) []A {
	ret := make([]A, 0, len(src))
	for _, val := range src {
		if v, ok := val[colKey]; ok {
			ret = append(ret, v)
		}
	}
	return ret
}
