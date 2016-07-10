'''
These are imported and called during the go tests
'''


import gosnake


def callee_none_none():
    return None


def callee_three_none(s, i, b):
    return None


def callee_none_one():
    return "higo"


def caller_none_none():
    gosnake.gocall()
    pass
