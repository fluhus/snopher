import ctypes

lib = ctypes.CDLL('./string.dll')
repeat = lib.repeat

repeat.argtypes = [
    ctypes.c_char_p,
    ctypes.c_int64,
    ctypes.c_char_p,
    ctypes.c_int64,
]
repeat.restype = ctypes.c_char_p

# Reusable output buffer.
buf_size = 1000
buf = ctypes.create_string_buffer(buf_size)

result = repeat(b'Badger', 4, buf, buf_size)  # type(result) = bytes
print('Badger * 4 =', result.decode())

result = repeat(b'Snake', 5, buf, buf_size)
print('Snake * 5 =', result.decode())
