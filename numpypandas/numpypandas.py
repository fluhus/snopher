import ctypes

import pandas

lib = ctypes.CDLL('./numpypandas.so')
increase = lib.increase

increase.argtypes = [
    ctypes.POINTER(ctypes.c_int64),
    ctypes.c_int64,
    ctypes.c_int64,
]

people = pandas.DataFrame(
    {
        'name': ['Alice', 'Bob', 'Charlie'],
        'age': [20, 30, 40],
    }
)

# First we check the type.
ages = people.age
if str(ages.dtypes) != 'int64':
    raise TypeError(f'Expected type int64, got {ages.dtypes}')

values = ages.values  # type=numpy.Array
ptr = values.ctypes.data_as(ctypes.POINTER(ctypes.c_int64))

print('Before')
print(people)

print('After')
increase(ptr, len(people), 5)
print(people)
