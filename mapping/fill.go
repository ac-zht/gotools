package mapping

// FillKeys 使用指定的键和值填充Map
func FillKeys[T comparable, A any](keys []T, val A) map[T]A {
    mp := make(map[T]A, len(keys))
    for _, v := range keys {
        mp[v] = val
    }
    return mp
}
