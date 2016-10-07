#! /usr/bin/python

import imp
import os


dir = os.path.dirname(__file__)
utils = imp.load_source('utils', os.path.join(dir, './utils.py'))
auth = imp.load_source('auth', os.path.join(dir, './auth/k8s.py'))

def do_test(t):
    name = t.__class__.__name__
    try:
        t.test()
        print ">>> %s: PASSED" % name
    except:
        print ">>> %s: FAILED" % name

if __name__ == "__main__":
    do_test(auth.Auth('incipit.machine:30002'))
