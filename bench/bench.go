package main

import "C"
import (
	"math/rand"
	"unsafe"
)

//export noop
func noop() {}

//export pi
func pi(n int64) float64 {
	result := 0.0
	sign := 1.0
	denom := 1.0
	for i := int64(0); i < n; i++ {
		result += sign * 4 / denom
		sign *= -1
		denom += 2
	}
	return result
}

//export shuffle
func shuffle(pnums *float64, n int64) {
	r := rand.New(rand.NewSource(0))
	nums := unsafe.Slice(pnums, n)
	r.Shuffle(int(n), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

//export dot
func dot(pa *float64, na int64, pb *float64, nb int64) float64 {
	a := unsafe.Slice(pa, na)
	b := unsafe.Slice(pb, nb)
	result := 0.0
	for i := int64(0); i < na; i++ {
		result += a[i] * b[i]
	}
	return result
}

func main() {}
