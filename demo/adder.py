'''
Demo python script integrating with gosnake
'''

import gosnake


def birthday(*args):
    # gosnake.gocall()
    return args


# dynamically call a method on gosnak
def callback(name, *args, **kwargs):
    print "PY: ", name, args, kwargs

    fromgo = gosnake.gocall(name, args)

    print "PY: go call returned with", fromgo

    # func = eval("gosnake." + name)
    # func(args, kwargs)
