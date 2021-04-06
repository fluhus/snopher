import ctypes


class UserInfo(ctypes.Structure):
    _fields_ = [
        ('name', ctypes.c_char_p),
        ('description', ctypes.c_char_p),
        ('name_length', ctypes.c_longlong),
    ]

    def __del__(self):
        del_user_info(self)


lib = ctypes.CDLL('user.dll')
get_user_info = lib.getUserInfo
get_user_info.argtypes = [ctypes.c_char_p]
get_user_info.restype = UserInfo
del_user_info = lib.delUserInfo
del_user_info.argtypes = [UserInfo]


def work_work():
    user1 = get_user_info('Alice'.encode())
    print('Name:', user1.name.decode())
    print('Description:', user1.description.decode())
    print('Name length:', user1.name_length)
    print('-----------')

    user2 = get_user_info('Bob'.encode())
    print('Name:', user2.name.decode())
    print('Description:', user2.description.decode())
    print('Name length:', user2.name_length)
    print('-----------')

    # Now user1 and user2 should get deleted.


work_work()
print('Did I remember to free my memory?')
