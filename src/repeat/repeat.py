import ctypes

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
