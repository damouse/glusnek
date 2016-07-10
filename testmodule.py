'''
These are imported and called during the go tests
'''


import gosnake

#
# Called from Go
#


def callee_none_none():
    return None


def callee_three_none(s, i, b):
    return None


def callee_none_one():
    return "higo"


#
# Calling into go
#
# Immediately call back into go
def reflect_call(name, *args):
    return gosnake.gocall(name, *args)
