'''
Demo python script integrating with gosnake
'''

import gosnake


def birthday(*args):
    print "PY: birthday- ", args
    # gosnake.gocall()

    return args


# dynamically call a method on gosnak
def callback(name, *args, **kwargs):
    func = eval("gosnake." + name)
    func(args, kwargs)
