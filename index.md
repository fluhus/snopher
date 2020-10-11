A tutorial and cheatsheet on calling Go code from Python using the ctypes library.

# Introduction

TODO: Add requirements, caveats.

# Hello World

Let's start with the bare minimum.

hello.go:

```go
package main

import "C"
import "fmt"

//export hello
func hello() {
    fmt.Println("Hello world!")
}

func main() {}
```

Build:

```
# Windows:
go build -o hello.dll -buildmode=c-shared hello.go
# Linux:
go build -o hello.so -buildmode=c-shared hello.go
```

hello.py:

```python
import ctypes

lib = ctypes.CDLL('./hello.dll')  # Or hello.so if on Linux.
hello = lib.hello

hello()
```

Then run:

```
> python hello.py
Hello world!
>
```

TODO: Add notes.

# Primitive Input and Output

TODO: Point out type mapping in .h file and architecture safety. Mention error handling.

add.go:

```go
//export add
func add(a, b int64) int64 {
    return a + b
}
```

add.py:

```python
lib = ctypes.CDLL('./add.dll')
add = lib.add

# Make python convert its values to C representation.
add.argtypes = [ctypes.c_longlong, ctypes.c_longlong]
add.restype = ctypes.c_longlong

print('10 + 15 =', add(10, 15))
```

Run:

```
> python add.py
10 + 15 = 25
>
```

TODO: Discuss argtypes and restype.

# Strings

TODO: Mention separate memory spaces.

repeat.go:

```go
//export repeat
func repeat(s *C.char, n int64) *C.char {
	var goString string = C.GoString(s) // Copy input to Go memory space.
	result := ""
	for i := int64(0); i < n; i++ {
		result += goString
	}
	return C.CString(result) // Copy result to C memory space.
}
```

repeat.py:

```python
lib = ctypes.CDLL('./repeat.dll')
repeat = lib.repeat

repeat.argtypes = [ctypes.c_char_p, ctypes.c_longlong]
repeat.restype = ctypes.c_char_p

result = repeat(b'Pizza', 4)  # type(result) = bytes
print('Pizza * 4 =', result.decode())
```

Run:

```
> python repeat.py
Pizza * 4 = PizzaPizzaPizzaPizza
>
```

# Arrays and Slices

# Structs

## Multiple Return Values

## Strings and Slices

# Benchmarks


