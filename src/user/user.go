package main

/*
#include <stdlib.h>
struct userInfo {
  char* info;
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Generates a data object for Python.
//
//export getUserInfo
func getUserInfo(cname *C.char) C.struct_userInfo {
	var result C.struct_userInfo
	name := C.GoString(cname)
	result.info = C.CString(
		fmt.Sprintf("User %q has %v letters in their name",
			name, len(name)))
	return result
}

// Deallocates a data object.
//
//export delUserInfo
func delUserInfo(info C.struct_userInfo) {
	// This print is only for educational purposes.
	fmt.Printf("Freeing user info: %s\n", C.GoString(info.info))
	C.free(unsafe.Pointer(info.info))
}

func main() {}
