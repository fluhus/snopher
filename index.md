A tutorial and cheatsheet on calling Go code from Python, using the ctypes library.

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

TODO: Point out type mapping in .h file and architecture safety.

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

TODO: Discuss argtypes and restype.

# Strings

# Arrays and Slices

# Structs

## Multiple Return Values

## Strings and Slices

# Benchmarks


