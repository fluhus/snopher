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

# Memory Spaces

**THIS SECTION IS IMPORTANT**

It is important to understand how memory spaces work in order to make efficient use
of Python and Go together.

In our case we have 3 memory spaces:

* **Python memory:** memory allocated by Python
* **Go memory:** memory allocated by Go's regular allocations
* **C memory:** memory allocated by Go using the "C" package

#### Python Memory

Limitations:

* If you **CONTINUE**

# Strings

TODO: Mention separate memory spaces.

repeat.go:

```go
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
```

repeat.py:

```python
lib = ctypes.CDLL('./repeat.dll')
repeat = lib.repeat

repeat.argtypes = [
    ctypes.c_char_p,
    ctypes.c_longlong,
    ctypes.c_char_p,
    ctypes.c_longlong,
]
repeat.restype = ctypes.c_char_p

# Reusable output buffer.
buf_size = 1000
buf = (ctypes.c_char * buf_size)(*([0] * buf_size))

result = repeat(b'Badger', 4, buf, buf_size)  # type(result) = bytes
print('Badger * 4 =', result.decode())

result = repeat(b'Snake', 5, buf, buf_size)
print('Snake * 5 =', result.decode())
```

Run:

```
> python repeat.py
Badger * 4 = BadgerBadgerBadgerBadger
Snake * 5 = SnakeSnakeSnakeSnakeSnake
>
```

# Arrays and Slices

# Structs

## Multiple Return Values

## Strings and Slices

# Benchmarks

A few comparisons to illustrate the potential benefit of using Go. All measurements include
the overhead of converting values to and from their C representation.

Tested on my personal desktop: Intel i5-6600K, 16GB RAM, Windows 10, Python 3.7.6, Go 1.14.

#### Pi

A simple comparison calculating the number Pi, to get a feeling of how much faster Go can be.

![pi](https://raw.githubusercontent.com/fluhus/snopher/master/pi.png)

#### Random Permutation

Comparing a more complex procedure. Notice how using Go can be faster than Python's builtins.

![shuffle](https://raw.githubusercontent.com/fluhus/snopher/master/shuffle.png)

#### Using `array` for Conversion

Comparing using the constructor recommended by the `ctypes` package documentation, to using
`array` for converting Python values to C values.

```python
# Using ctypes
cvals = (ctypes.c_double * n)(*nums)

# Using array
arr = array('d', nums)
cvals = (ctypes.c_double * n).from_buffer(arr)
```

![list](https://raw.githubusercontent.com/fluhus/snopher/master/list.png)




