import ctypes

lib = ctypes.CDLL('./add.dll')
add = lib.add

# Make python convert its values to C representation.
add.argtypes = [ctypes.c_longlong, ctypes.c_longlong]
add.restype = ctypes.c_longlong

print('10 + 15 =', add(10, 15))
