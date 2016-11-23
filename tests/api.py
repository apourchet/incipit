#! /usr/bin/python

import imp
import os
import sys
import traceback
from termcolor import colored

dir = os.path.dirname(__file__)
utils = imp.load_source('utils', os.path.join(dir, './utils.py'))
auth = imp.load_source('auth', os.path.join(dir, './auth/api.py'))

START = colored("START", "white")
PASSED = colored("PASSED", "green")
FAILED = colored("FAILED", "red")

def do_test(t):
    name = t.__class__.__name__
    try:
        print ">>> %s: %s" % (name, START)
        t.test()
        print ">>> %s: %s" % (name, PASSED)
    except Exception as e:
        _, _, tb = sys.exc_info()
        info = traceback.extract_tb(tb)
        filename, line, func, text = info[-1]
        print "File '%s'; Line %d; Statement '%s'" % (filename, line, text)
        print ">>> %s: %s" % (name, FAILED)

if __name__ == "__main__":
    do_test(auth.Auth('incipit.machine:30002'))
