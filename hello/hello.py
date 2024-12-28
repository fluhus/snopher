import ctypes

lib = ctypes.CDLL('./hello.so')  # Or hello.dll if on Windows.
hello = lib.hello

hello()
