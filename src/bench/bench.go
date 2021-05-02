package main

import "C"
import (
	"math/rand"
	"unsafe"
)

//export noop
func noop() {}

//export pi
func pi(n int) float64 {
	result := 0.0
	sign := 1.0
	denom := 1.0
	for i := 0; i < n; i++ {
		result += sign * 4 / denom
		sign *= -1
		denom += 2
	}
	return result
}

//export shuffle
func shuffle(pnums *float64, n int) {
	r := rand.New(rand.NewSource(0))
	nums := (*[1 << 30]float64)(unsafe.Pointer(pnums))[:n:n]
	r.Shuffle(n, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
}

//export dot
func dot(pa *float64, na int, pb *float64, nb int) float64 {
	a := (*[1 << 30]float64)(unsafe.Pointer(pa))[:na:na]
	b := (*[1 << 30]float64)(unsafe.Pointer(pb))[:nb:nb]
	result := 0.0
	for i := 0; i < na; i++ {
		result += a[i] * b[i]
	}
	return result
}

func main() {}
