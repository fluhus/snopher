import ctypes
from typing import List


def to_c_str_array(strs: List[str]):
    ptr = (ctypes.c_char_p * (len(strs) + 1))()
    ptr[:-1] = [s.encode() for s in strs]
    ptr[-1] = None  # Terminate with null.
    return ptr


lib = ctypes.CDLL('./join.dll')
join = lib.join
join.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_char_p]

words = ['yab', 'dab', 'doo!']
sep = 'a'

join(to_c_str_array(words), sep.encode())
