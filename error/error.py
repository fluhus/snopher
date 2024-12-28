import ctypes


class Error(ctypes.Structure):
    _fields_ = [('err', ctypes.c_char_p)]

    def __del__(self):
        # We can call del_error with a None err, but this way we can avoid
        # the call overhead when it's not necessary.
        if self.err is not None:
            del_error(self)

    def raise_if_err(self):
        if self.err is not None:
            raise IOError(self.err.decode())


class EvenResult(ctypes.Structure):
    # Multiple return values can be grabbed in a struct.
    _fields_ = [
        ('result', ctypes.c_bool),
        ('err', Error),
    ]


lib = ctypes.CDLL('./error.dll')

del_error = lib.delError
del_error.argtypes = [Error]

even = lib.even
even.argtypes = [ctypes.c_int64]
even.restype = EvenResult

for i in (0, 1, 2, 3, -5):
    e = even(i)
    e.err.raise_if_err()
    print(i, 'even:', e.result)
