A tutorial and cheatsheet on calling Go code from Python, using the ctypes library.

# Introduction

TODO: Add requirements.

## Caveats

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

# Strings

# Arrays and Slices

# Structs

## Multiple Return Values

## Strings and Slices

# Benchmarks


