package main

import "C"
import "unsafe"

//export increase
func increase(numsPtr *int64, n int64, a int64) {
	nums := (*[1 << 30]int64)(unsafe.Pointer(numsPtr))[:n:n]
	for i := range nums {
		nums[i] += a
	}
}

func main() {}
