package slice

func Diff[T comparable](s1 []T, s2 []T) []T {
    m1 := Map[T](s1)
    for _, v := range s2 {
        delete(m1, v)
    }
}
