import ctypes
from array import array

# A function that receives an array type string and a size,
# and returns a pointer.
alloc_f = ctypes.CFUNCTYPE(ctypes.c_void_p, ctypes.c_char_p, ctypes.c_int64)

arrays: list[array] = []


@alloc_f
def my_alloc(typecode, size):
    arr = array(typecode.decode(), (0 for _ in range(size)))
    arrays.append(arr)
    return arr.buffer_info()[0]


lib = ctypes.CDLL('./alloc.so')
sqrts = lib.sqrts
sqrts.argtypes = [alloc_f, ctypes.c_int64]

sqrts(my_alloc, 5)  # Populates arrays with the result.

print(arrays[0].tobytes().decode(), arrays[1].tolist())
