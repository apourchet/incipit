#! /usr/bin/python

from random import randint
import sys
from subprocess import Popen, PIPE

def newport():
    return randint(31000, 32000)

def expose(servicename, port=newport()):
    cmd = 'kubectl patch service %s --type=\'json\' ' % servicename
    cmd += '-p=\'[{"op": "replace", "path": "/spec/type", "value": "NodePort"},'
    cmd += '{"op": "replace", "path": "/spec/ports/0/nodePort", "value": "%d"}]\'' % port
    output, err = Popen(cmd, shell=True, executable='/bin/bash', stdout=PIPE, stderr=PIPE).communicate()
    if err != None and err != "":
        return None
    return port

def hide(servicename):
    cmd = 'kubectl patch service %s --type=\'json\' ' % servicename
    cmd += '-p=\'[{"op": "replace", "path": "/spec/type", "value": "ClusterIP"},'
    cmd += '{"op": "remove", "path": "/spec/ports/0/nodePort"}]\''
    return Popen(cmd, shell=True, executable='/bin/bash', stdout=PIPE, stderr=PIPE).communicate()

if __name__ == "__main__":
    if len(sys.argv) != 3 or (sys.argv[2] != 'expose' and sys.argv[2] != 'hide'):
        print "Usage: kubexpose.rb <servicename> <expose/hide>"
        exit(0)

    servicename = sys.argv[1]
    action = sys.argv[2]

    if action == 'expose':
        port = expose(servicename)
        if port != None:
            print port
        else:
            print "Error exposing"
            exit(1)
    else:
        hide(servicename)
