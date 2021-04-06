package main

import "C"
import "unsafe"

// Returns the squares of the input numbers.
//
//export squares
func squares(numsPtr *float64, outPtr *float64, n int64) {
	// The way to wrap a pointer with a Go slice.
	nums := (*[1 << 30]float64)(unsafe.Pointer(numsPtr))[:n:n]
	out := (*[1 << 30]float64)(unsafe.Pointer(outPtr))[:n:n]

	for i, x := range nums {
		out[i] = x * x
	}
}

func main() {}
