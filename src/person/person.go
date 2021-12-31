package main

/*
struct person {
  char* firstName;
  char* lastName;
  char* fullName;
  long long fullNameLen;
};
*/
import "C"
import (
	"bytes"
	"unsafe"
)

//export fill
func fill(p *C.struct_person) {
	buf := bytes.NewBuffer(unsafe.Slice((*byte)(unsafe.Pointer(p.fullName)),
		p.fullNameLen)[:0])
	first := C.GoString(p.firstName)
	last := C.GoString(p.lastName)
	buf.WriteString(first + " " + last)
	buf.WriteByte(0)
}

func main() {}
