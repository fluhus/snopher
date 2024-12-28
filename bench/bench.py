import ctypes
import random
from array import array
from time import monotonic
from typing import Callable


def humanize_seconds(t):
    if t < 0.001:
        return f'{t*1000000:.1f}us'
    elif t < 1:
        return f'{t*1000:.1f}ms'
    else:
        return f'{t:.1f}s'


def benchmark(f: Callable, name: str):
    n = 0
    diff = 0
    while diff < 1.5:
        n = (n + 1) * 3 // 2
        t = monotonic()
        for _ in range(n):
            v = f()
        diff = monotonic() - t
    print(f'{name:15}{humanize_seconds(diff/n)} / iteration')
    return v


lib = ctypes.CDLL('./bench.dll')


def benchmark_noop():
    print('*** No-op ***')
    noop = lib.noop
    benchmark(noop, 'No-op')


def benchmark_pi():
    print('*** Pi ***')

    pi = lib.pi
    pi.argtypes = [ctypes.c_int64]
    pi.restype = ctypes.c_double

    def pypi(n):
        p = 0
        sign = 1
        for i in range(n):
            p += sign * 4 / (i * 2 + 1)
            sign *= -1
        return p

    n = 10000000
    pi1 = benchmark(lambda: pi(n), 'Go pi')
    pi2 = benchmark(lambda: pypi(n), 'Py pi')

    print('Go pi =', pi1)
    print('Py pi =', pi2)


def benchmark_list_conversion():
    print('*** List Convertions ***')

    def to_array(x):
        y = array('d', x)
        (ctypes.c_double * len(y)).from_buffer(y)
        return y

    n = 10000000
    nums = benchmark(lambda: list(range(n)), 'Alloc')
    buf = benchmark(lambda: (ctypes.c_double * n)(*nums), 'CTypes-to')
    arr = benchmark(lambda: to_array(nums), 'Array-to')
    benchmark(lambda: list(buf), 'Ctypes-from')
    benchmark(lambda: list(arr), 'Array-from')


def benchmark_shuffle():
    print('*** Shuffle ***')
    shuffle = lib.shuffle
    shuffle.argtypes = [ctypes.POINTER(ctypes.c_double), ctypes.c_int64]

    def go_shuffle(x):
        arr = array('d', nums)
        buf = (ctypes.c_double * n).from_buffer(arr)
        shuffle(buf, len(x))
        list(arr)

    n = 1000000
    nums = benchmark(lambda: list(range(n)), 'Alloc')
    benchmark(lambda: go_shuffle(nums), 'Go')
    benchmark(lambda: random.shuffle(nums), 'Random')

    # Importing numpy.random makes ctypes.CDLL not find the dll. :-\
    from numpy import random as nprandom

    benchmark(lambda: nprandom.shuffle(nums), 'Numpy')


def benchmark_dot():
    print('*** Dot ***')
    dot = lib.dot
    dot.argtypes = [
        ctypes.POINTER(ctypes.c_double),
        ctypes.c_int64,
        ctypes.POINTER(ctypes.c_double),
        ctypes.c_int64,
    ]
    dot.restype = ctypes.c_double

    import numpy

    n = 100000000
    arr1, arr2 = benchmark(
        lambda: numpy.ones([2, n], dtype=numpy.float64),
        'Alloc',
    )

    p1 = arr1.ctypes.data_as(ctypes.POINTER(ctypes.c_double))
    p2 = arr2.ctypes.data_as(ctypes.POINTER(ctypes.c_double))
    benchmark(lambda: dot(p1, len(arr1), p2, len(arr2)), 'Go')
    benchmark(lambda: arr1.dot(arr2), 'Numpy')


benchmark_noop()
print()
benchmark_pi()
print()
benchmark_list_conversion()
print()
benchmark_shuffle()
print()
benchmark_dot()
