package util

import (
	"math/rand"
	"time"
)

// Random 随机数
func Random(args ...int) int {
	var (
		min, max = crandom(args...)
	)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// RandomMutil 返回多个随机数
func RandomMutil(count int, args ...int) []int {
	min, max := crandom(args...)
	if max-min < count {
		count = max - min
	}
	arr := make([]int, 0, count)
	newarr := make([]int, 0, count)
	for i := min; i <= max; i++ {
		arr = append(arr, i)
	}
	for i := 0; i < count; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(arr))
		newarr = append(newarr, arr[r])
		arr = append(arr[:r], arr[r+1:]...)
	}
	return newarr
}

// Random 随机数 =>
func crandom(args ...int) (min, max int) {
	if len(args) == 2 {
		min, max = args[0], args[1]
	} else if len(args) == 1 {
		max = args[0]
	}
	return
}
