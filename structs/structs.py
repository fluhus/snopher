import ctypes


class Person(ctypes.Structure):
    _fields_ = [
        ('first_name', ctypes.c_char_p),
        ('last_name', ctypes.c_char_p),
        ('full_name', ctypes.c_char_p),
        ('full_name_len', ctypes.c_int64),
    ]


lib = ctypes.CDLL('./structs.so')

fill = lib.fill
fill.argtypes = [ctypes.POINTER(Person)]

buf_size = 1000
buf = ctypes.create_string_buffer(buf_size)
person = Person(b'John', b'Galt', buf.value, len(buf))
fill(ctypes.pointer(person))

print(person.full_name.decode())
