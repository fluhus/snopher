package main

/*
#include <stdlib.h>
typedef struct {
	char* err;
} error;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func newError(s string, args ...interface{}) C.error {
	if s == "" {
		// Equivalent to a nil Go error.
		return C.error{}
	}
	msg := fmt.Sprintf(s, args...)
	return C.error{C.CString(msg)}
}

//export delError
func delError(err C.error) {
	if err.err == nil {
		return
	}
	C.free(unsafe.Pointer(err.err))
}

// Checks if a non-negative number is even.
//
//export even
func even(i int64) (bool, C.error) {
	if i < 0 {
		return false, newError("%v is negative, want at least 0", i)
	}
	return i%2 == 0, newError("")
}

func main() {}
