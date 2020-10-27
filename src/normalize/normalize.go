package main

import "C"
import "unsafe"

// Returns the input numbers minus their mean.
//
//export normalize
func normalize(numsPtr *float64, n int64, outPtr *float64) {
	// The way to wrap a pointer with a Go slice.
	nums := (*[1 << 30]float64)(unsafe.Pointer(numsPtr))[:n:n]
	out := (*[1 << 30]float64)(unsafe.Pointer(outPtr))[:n:n]

	// Calculate mean.
	mean := 0.0
	for _, num := range nums {
		mean += num
	}
	mean /= float64(n)

	// Assign output.
	for i := range nums {
		out[i] = nums[i] - mean
	}
}

func main() {}
