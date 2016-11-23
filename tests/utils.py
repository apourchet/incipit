#! /usr/bin/python

import requests
import time
import imp
import os
from requests.packages.urllib3.exceptions import InsecureRequestWarning
from requests.adapters import HTTPAdapter
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

class Session(requests.Session):
    timeout = 1
    retries = 10

    def set_timeout(self, timeout):
        self.timeout = timeout

    def set_retries(self, retries):
        self.retries = retries

    def do_retry(self, fn):
        count = 0
        while count < self.retries:
            count += 1
            try:
                resp = fn()
                return resp
            except (requests.exceptions.ReadTimeout, requests.exceptions.ConnectionError) as e:
                time.sleep(1)

    def post(self, *args, **kwargs):
        kwargs['timeout'] = self.timeout
        return self.do_retry(lambda: super(Session, self).post(*args, **kwargs))

    def get(self, *args, **kwargs):
        kwargs['timeout'] = self.timeout
        return self.do_retry(lambda: super(Session, self).get(*args, **kwargs))

def newSession():
    session = Session()
    session.verify = False
    session.headers.update({
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    })
    return session

def setSessionToken(session, token):
    session.headers.update({ 'Authorization': 'Bearer %s' % token })

def doImport(name, path):
    dir = os.path.dirname(__file__)
    return imp.load_source(name, os.path.join(dir, path))

