package slice

type equalFunc[T comparable] func(src, dst T) bool
