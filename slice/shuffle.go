package slice

import (
    "github.com/zht-account/gotools/random"
    "math/rand"
    "reflect"
    "time"
)

// Shuffle 打乱切片
func Shuffle[T any](src []T) {
    typ := reflect.TypeOf(src)
    if typ.Kind() != reflect.Slice {
        return
    }
    if len(src) <= 1 {
        return
    }
    rand.Seed(time.Now().UnixNano())
    srcSwap := reflect.Swapper(src)
    var j int
    for i := 0; i < len(src); i++ {
        j = random.RandInt(0, len(src)-1)
        srcSwap(i, j)
    }
}
