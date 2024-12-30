package main

/*
#include <stdlib.h>
#include <stdint.h>

typedef void* (*alloc_f)(char* t, int64_t n);
static void* call_alloc_f(alloc_f f, char* t, int64_t n) {return f(t,n);}
*/
import "C"
import (
	"fmt"
	"math"
	"unsafe"
)

func main() {}

//export sqrts
func sqrts(alloc C.alloc_f, n int64) {
	allocString(alloc, fmt.Sprintf("Square roots up to %d:", n))
	floats := allocFloats(alloc, int(n))
	for i := range floats {
		floats[i] = math.Sqrt(float64(i + 1))
	}
}

func allocFloats(alloc C.alloc_f, n int) []float64 {
	return allocSlice[float64](alloc, n, "d")
}

func allocInts(alloc C.alloc_f, n int) []int64 {
	return allocSlice[int64](alloc, n, "q")
}

func allocBytes(alloc C.alloc_f, n int) []byte {
	return allocSlice[byte](alloc, n, "B")
}

func allocString(alloc C.alloc_f, s string) {
	b := allocBytes(alloc, len(s))
	copy(b, s)
}

// Calls the Python alloc callback and returns the allocated buffer
// as a slice.
func allocSlice[T any](alloc C.alloc_f, n int, typeCode string) []T {
	t := C.CString(typeCode)                      // Make a c-string type code.
	ptr := C.call_alloc_f(alloc, t, C.int64_t(n)) // Allocate the buffer.
	C.free(unsafe.Pointer(t))                     // Release c-string.
	return unsafe.Slice((*T)(ptr), n)             // Wrap with a go-slice.
}
