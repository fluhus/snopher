package main

import "C"
import "unsafe"

// Returns the squares of the input numbers.
//
//export squares
func squares(numsPtr *float64, outPtr *float64, n int64) {
	// Wrap the pointers with Go slices (pointing to the same data).
	nums := unsafe.Slice(numsPtr, n)
	out := unsafe.Slice(outPtr, n)

	// If using Go < 1.17.
	// nums := (*[1 << 30]float64)(unsafe.Pointer(numsPtr))[:n:n]
	// out := (*[1 << 30]float64)(unsafe.Pointer(outPtr))[:n:n]

	for i, x := range nums {
		out[i] = x * x
	}
}

func main() {}
