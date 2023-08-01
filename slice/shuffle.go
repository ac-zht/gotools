package slice

import (
	"math/rand"
	"reflect"
	"time"
)

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
		j = rand.Intn(len(src))
		srcSwap(i, j)
	}
}
