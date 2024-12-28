import ctypes

lib = ctypes.CDLL('./hello.dll')  # Or hello.so if on Linux.
hello = lib.hello

hello()
