**UNDER CONSTRUCTION**

This is a tutorial and cheatsheet on calling Go code from Python using the
ctypes library.

It is based on the experience I accumulated while working on my
research projects. It may contain suboptimal practices or even errors. Feedback
and ideas are welcome!

Snake + Gopher = <3

Author: Amit Lavon

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

Congrats!

Let's break it down:

1. The Go code uses regular Go logic, but exports its function for external
   use with the `//export` directive.
2. Building with `-buildmode=c-shared` creates a C-style shared library.
3. Python loads the shared library and accesses the exported function.

# Primitive Input and Output

Here we introduce some basic arguments and return values. Let's start with an
example.

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

To pass input and receive output to/from a Go function, we need to use ctypes's
`argtypes` and `restype` attributes. They do 2 things:

1. `argtypes` guards the function by checking the arguments before calling the
   library code.
2. Using these attributes tells Python how to convert the input Python values
   to ctypes values, and how to convert the output back to a Python value.

You can find the mapping between C types and Go types in the generated `.h`
file after you compile your Go code with `-buildmode=c-shared`.

**DANGER: Some types change sizes on different architectures. It is generally
safer to use sized types (`int64`) here than unsized types (`int`).**

# Arrays and Slices

We are now entering the dangerous zone of unprotected memory access. While
Python and Go are generally memory safe, working with raw pointers might end
up in buffer overflows and memory leaks.

**Make sure to read this section through in order to avoid bad things.**

normalize.go:

```go
// Returns the input numbers minus their mean.
//
//export normalize
func normalize(numsPtr *float64, outPtr *float64, n int64) {
	// The way to wrap a pointer with a Go slice.
	nums := (*[1 << 30]float64)(unsafe.Pointer(numsPtr))[:n:n]
	out := (*[1 << 30]float64)(unsafe.Pointer(outPtr))[:n:n]

	// Calculate mean.
	mean := 0.0
	for _, num := range nums {
		mean += num
	}
	mean /= float64(n)

	// Assign output.
	for i := range nums {
		out[i] = nums[i] - mean
	}
}
```

normalize.py:

```python
lib = ctypes.CDLL('./normalize.dll')
normalize = lib.normalize

normalize.argtypes = [
    ctypes.POINTER(ctypes.c_double),
    ctypes.POINTER(ctypes.c_double),
    ctypes.c_longlong,
]

# Building buffers from arrays is more efficient than
# (ctypes.c_double * 3)(*[1, 2, 3])
nums = array('d', [1, 2, 3])
nums_ptr = (ctypes.c_double * len(nums)).from_buffer(nums)
out = array('d', (0 for _ in range(len(nums))))
out_ptr = (ctypes.c_double * len(out)).from_buffer(out)

normalize(nums_ptr, out_ptr, len(nums))
print('nums:', list(nums))
print('out:', list(out))
```

Run:

```
> python normalize.py
nums: [1.0, 2.0, 3.0]
out: [-1.0, 0.0, 1.0]
>
```

In order to work with lists we need to convert them into C arrays. The way to
do it is to create an array using `(ctypes.my_type * my_length)(1, 2, 3 ...)`.
The faster way to do it is to use the `array` library as demonstrated above.
See the benchmarks section for more info about their performance.

In Go you can cast the C-like pointer to a slice (that still points to Python's
memory), see the demonstration above. This way you can utilize Go's syntax while
working with Python buffers.

Output lists are the tricky part. You cannot return a Go pointer when using CGo,
that will result in an error. Instead you can allocate a C pointer from Go using
`C.malloc()` and return that. However, that pointer is not garbage-collected
so unless you implement a way to deallocate those, you will have a memory leak.

The possibly safest way to go about output arrays is to pre-allocate them in
Python and pass them as arguments to the function. Notice that you need to keep
their references in your python code until you are done running Go code on them,
to keep them from getting garbage-collected.

#### Summary of Dangers

* Returning Go pointers to Python. **Error.**
* Returning C pointers from Go to Python without explicit deallocation.
  **Memory leak.**
* Losing the `ctypes` reference while Go code is still running (like when taking
  `ctypes.addressof` and dumping the pointer object). **Possible segmentation
  fault.**

# Strings

Strings work pretty much like arrays in terms of memory management, so
everything related to arrays applies here too. Below I discuss some convenience
techniques and some pitfalls.

repeat.go:

```go
//export repeat
func repeat(s *C.char, n int64, out *C.char, outN int64) *C.char {
	// Create a Go buffer around the output buffer.
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

Strings are passed by converting a Python string to a bytes object (typically
by calling `encode()`), then to a C-pointer and then to a Go-string.

Using `ctypes.c_char_p` in argtypes makes Python expect a bytes object and
convert it to a C `*char`. In restype, it converts the returned `*char` to a
bytes object.

In Go you can convert a `*char` to a Go string using `C.GoString`. This copies
the data and creates a new string managed by Go in terms of garbage collection.

To create a `*char` as a return value, you can call `C.CString`. However, the
pointer gets lost unless you keep a reference to it in Go, and then you have a
memory leak.

My recommended way to return a string is to create an output buffer in Python
(possibly reusable), pass it to Go and then wrap it in a `bytes.Buffer`. That
should make generating the output string safe and efficient. Don't forget the
null terminator!

Go can return the given output pointer, and Python will automatically make a
bytes object out of it. See the demonstration above.

#### Summary of Dangers

* Returning a `C.CString` without keeping the reference for future deallocation.
  **Memory leak.**
* Not appending a null terminator at the end of the output string.
  **Buffer overflow when converting to Python object.**
* Not checking output buffer size in Go. **Buffer overflow or truncated
  output.**

# Numpy and Pandas

Numpy provides access to its underlying buffers using the
`.ctypes.data_as(ctypes.whatever)` syntax. With pandas you can use the `.values`
attribute to get the underlying numpy array, and then use numpy's syntax to get
the actual pointer. This way you can change the array/table in place.

table.go:

```go
//export increase
func increase(numsPtr *int64, n int64, a int64) {
	nums := (*[1 << 30]int64)(unsafe.Pointer(numsPtr))[:n:n]
	for i := range nums {
		nums[i] += a
	}
}
```

table.py:

```python
lib = ctypes.CDLL('./table.dll')
increase = lib.increase

increase.argtypes = [
    ctypes.POINTER(ctypes.c_longlong),
    ctypes.c_longlong,
    ctypes.c_longlong,
]

people = pandas.DataFrame(
    {
        'name': ['Alice', 'Bob', 'Charlie'],
        'age': [20, 30, 40],
    }
)

# First we check the type.
ages = people.age
if str(ages.dtypes) != 'int64':
    raise TypeError(f'Expected type int64, got {ages.dtypes}')

values = ages.values  # type=numpy.Array
ptr = values.ctypes.data_as(ctypes.POINTER(ctypes.c_longlong))

print('Before')
print(people)

print('After')
increase(ptr, len(people), 5)
print(people)
```

Run:

```
> python table.py
Before
      name  age
0    Alice   20
1      Bob   30
2  Charlie   40
After
      name  age
0    Alice   25
1      Bob   35
2  Charlie   45
>
```

# Structs

## Multiple Return Values

## Strings and Slices

# Performance Tips

Mention reusing buffers

Mention using `array`

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




