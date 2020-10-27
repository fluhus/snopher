package main

import "C"
import (
	"bytes"
	"unsafe"
)

//export repeat
func repeat(s *C.char, n int64, out *C.char, outN int64) *C.char {
	// Create a Go buffer around output buffer.
	outBytes := (*[1 << 30]byte)(unsafe.Pointer(out))[:0:outN]
	buf := bytes.NewBuffer(outBytes)

	var goString string = C.GoString(s) // Copy input to Go memory space.
	for i := int64(0); i < n; i++ {
		buf.WriteString(goString)
	}
	buf.WriteByte(0) // Null terminator, important!
	return out
}

func main() {}
