import ctypes

lib = ctypes.CDLL('./primitive.so')
add = lib.add

# Make python convert its values to C representation.
add.argtypes = [ctypes.c_int64, ctypes.c_int64]
add.restype = ctypes.c_int64

print('10 + 15 =', add(10, 15))
