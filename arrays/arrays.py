import ctypes
from array import array

lib = ctypes.CDLL('./arrays.so')
squares = lib.squares

squares.argtypes = [
    ctypes.POINTER(ctypes.c_double),
    ctypes.POINTER(ctypes.c_double),
    ctypes.c_int64,
]

# Building buffers from arrays is more efficient than
# (ctypes.c_double * 3)(*[1, 2, 3])
nums = array('d', [1, 2, 3])
nums_ptr = (ctypes.c_double * len(nums)).from_buffer(nums)
out = array('d', [0, 0, 0])
out_ptr = (ctypes.c_double * len(out)).from_buffer(out)

squares(nums_ptr, out_ptr, len(nums))
print('nums:', list(nums), 'out:', list(out))
