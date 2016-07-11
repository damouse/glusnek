'''
Demo python script integrating with gosnake
'''

import gosnake

import json
import threading
import datetime

counter = 0


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
    global counter
    counter += 1

    ret = gosnake.gocall("Hello!")
    # print "PY: args: ", a, " goret: ", ret

    # Holy crap this works
    # print "Starting a sleep"
    # import time
    # time.sleep(2)  # delays for 5 seconds

    return json.dumps([threading.currentThread().ident, counter] + list(a))


#
# Sandbox
#
def sandy(*args):
    print "PY: ", args
