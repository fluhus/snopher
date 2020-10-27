import random
import time
from contextlib import contextmanager
import ctypes
from array import array


@contextmanager
def timer(prefix=None):
    t = time.monotonic()
    yield
    t = time.monotonic() - t
    if prefix is not None:
        print('{:15}'.format(prefix), end='')
    if t < 0.001:
        print('{:.1f}us'.format(t * 1000000))
    elif t < 1:
        print('{:.1f}ms'.format(t * 1000))
    else:
        print('{:.1f}s'.format(t))


lib = ctypes.CDLL('./bench.dll')


def benchmark_pi():
    print('*** Pi ***')

    pi = lib.pi
    pi.argtypes = [ctypes.c_longlong]
    pi.restype = ctypes.c_double

    n = 10000000

    with timer('Go'):
        gopy = pi(n)

    with timer('Python'):
        pypi = sum((-1) ** i * 4 / (i * 2 + 1) for i in range(n))

    print('Go pi =', gopy)
    print('Py pi =', pypi)


def benchmark_list_conversion():
    print('*** List Convertions ***')

    n = 10000000
    with timer('Alloc'):
        nums = list(range(n))

    with timer('Ctypes-to'):
        buf = (ctypes.c_double * n)(*nums)

    with timer('Array-to'):
        arr = array('d', nums)
        (ctypes.c_double * n).from_buffer(arr)

    with timer('Ctypes-from'):
        list(buf)

    with timer('Array-from'):
        list(arr)


def benchmark_shuffle():
    print('*** Shuffle ***')
    shuffle = lib.shuffle
    shuffle.argtypes = [ctypes.POINTER(ctypes.c_double), ctypes.c_longlong]

    n = 10000000
    with timer('Alloc'):
        nums = list(range(n))

    with timer('Go'):
        arr = array('d', nums)
        buf = (ctypes.c_double * n).from_buffer(arr)
        shuffle(buf, n)
        list(arr)

    with timer('Random'):
        random.shuffle(nums)

    # Importing numpy.random make ctypes.CDLL not find the dll. :-\
    from numpy import random as nprandom

    with timer('Numpy'):
        nprandom.shuffle(nums)


def benchmark_dot():
    print('*** Dot ***')
    dot = lib.dot
    dot.argtypes = [
        ctypes.POINTER(ctypes.c_double),
        ctypes.c_longlong,
        ctypes.POINTER(ctypes.c_double),
        ctypes.c_longlong,
    ]
    dot.restype = ctypes.c_double

    import numpy

    n = 100000000
    t = 10

    with timer('Alloc'):
        arr1 = numpy.ndarray([n], dtype=numpy.float64)
        arr2 = numpy.ndarray([n])
        arr1[:] = 1
        arr2[:] = 1
    # print(arr1, arr2)

    with timer('Go'):
        for _ in range(t):
            p1 = arr1.ctypes.data_as(ctypes.POINTER(ctypes.c_double))
            p2 = arr2.ctypes.data_as(ctypes.POINTER(ctypes.c_double))
            dot(p1, len(arr1), p2, len(arr2))
    
    with timer('Numpy'):
        for _ in range(t):
            arr1.dot(arr2)


# benchmark_pi()
# print()
# benchmark_list_conversion()
# print()
benchmark_shuffle()
# print()
benchmark_dot()
