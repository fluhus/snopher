package main

/*
#include <stdlib.h>
struct userInfo {
  char* name;
  char* description;
  long long nameLength;
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
	// Create a copy to give it the same lifetime as the rest of the object.
	result.name = C.CString(name)
	result.description = C.CString(
		fmt.Sprintf("User %q has %v letters in their name",
			name, len(name)))
	result.nameLength = C.longlong(len(name))
	return result
}

// Deallocates a data object.
//
//export delUserInfo
func delUserInfo(info C.struct_userInfo) {
	fmt.Printf("Freeing user %q\n", C.GoString(info.name))
	C.free(unsafe.Pointer(info.name))
	C.free(unsafe.Pointer(info.description))
}

func main() {}
