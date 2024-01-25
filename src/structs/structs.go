package main

/*
#include <stdint.h>
typedef struct person {
  char* firstName;
  char* lastName;
  char* fullName;
  int64_t fullNameLen;
} person;
*/
import "C"
import (
	"bytes"
	"unsafe"
)

//export fill
func fill(p *C.person) {
	buf := bytes.NewBuffer(unsafe.Slice((*byte)(unsafe.Pointer(p.fullName)),
		p.fullNameLen)[:0])
	first := C.GoString(p.firstName)
	last := C.GoString(p.lastName)
	buf.WriteString(first + " " + last)
	buf.WriteByte(0)
}

func main() {}
