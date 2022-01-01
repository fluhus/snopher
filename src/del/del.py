import ctypes


class UserInfo(ctypes.Structure):
    _fields_ = [('info', ctypes.c_char_p)]

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
    print('Info:', user1.info.decode())
    print('-----------')

    user2 = get_user_info('Bob'.encode())
    print('Info:', user2.info.decode())
    print('-----------')

    # Now user1 and user2 should get deleted.


work_work()
print('Did I remember to free my memory?')
