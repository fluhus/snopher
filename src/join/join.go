package main

import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

func goStrings(cstrs **C.char) []string {
	var result []string
	slice := (*[1 << 30]*C.char)(unsafe.Pointer(cstrs))[: 1<<30 : 1<<30]
	for i := 0; slice[i] != nil; i++ {
		result = append(result, C.GoString(slice[i]))
	}
	return result
}

//export join
func join(strs **C.char, sep *C.char) {
	goStrs := goStrings(strs)
	goSep := C.GoString(sep)
	fmt.Println(strings.Join(goStrs, goSep))
}

func main() {}
