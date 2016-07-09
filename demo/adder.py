'''
Demo python script integrating with gosnake
'''

import gosnake

import json
import threading
import datetime


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


def run(*a):
    # print "PY: run: ", a
    return json.dumps([threading.currentThread().ident] + list(a))
