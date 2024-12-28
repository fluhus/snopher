package main

import "C"
import (
	"bytes"
	"unsafe"
)

//export repeat
func repeat(s *C.char, n int64, out *byte, outN int64) *byte {
	// Create a Go buffer around the output buffer.
	outBytes := unsafe.Slice(out, outN)[:0]
	buf := bytes.NewBuffer(outBytes)

	var goString string = C.GoString(s) // Copy input to Go memory space.
	for range n {
		buf.WriteString(goString)
	}
	buf.WriteByte(0) // Null terminator, important!
	return out
}

func main() {}
