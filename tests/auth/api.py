#! /usr/bin/python

import time
import random
import imp
import os
import requests
from requests.packages.urllib3.exceptions import InsecureRequestWarning
requests.packages.urllib3.disable_warnings(InsecureRequestWarning)

dir = os.path.dirname(__file__)
utils = imp.load_source('utils', os.path.join(dir, '../utils.py'))

class Auth():
    def __init__(self, ing_addr):
        self.ing_addr = ing_addr

    def register(self, session, username, password):
        url = 'https://%s/api/v1/auth/register' % self.ing_addr
        return session.post(url, json={"key": username, "pass": password})

    def login(self, session, username, password):
        url = 'https://%s/api/v1/auth/login' % self.ing_addr
        url += '?key=%s&pass=%s' % (username, password)
        resp = session.get(url)
        if 'token' in resp.json():
            token = resp.json()['token']
            utils.setSessionToken(session, token)
        return resp

    def userExists(self, session, username):
        url = 'https://%s/api/v1/auth/userexists' % self.ing_addr
        url += '?key=%s' % (username)
        return session.get(url)

    def logout(self, session):
        url = 'https://%s/api/v1/auth/logout' % self.ing_addr
        return session.post(url, json={})

    def deregister(self, session):
        url = 'https://%s/api/v1/auth/deregister' % self.ing_addr
        return session.post(url, json={})

    def test(self):
        session = utils.newSession()
        username = "auth-test-username-%d" % (random.randint(0, 10000))
        resp = self.register(session, username, "pass")
        assert resp.status_code == 200

        resp = self.login(session, username, "pass")
        assert resp.status_code == 200
        assert resp.json()['token'] != ""

        resp = self.userExists(session, username)
        assert resp.status_code == 200
        assert resp.json()['found'] == True

        resp = self.logout(session)
        assert resp.status_code == 200

        resp = self.deregister(session)
        assert resp.status_code == 401

        resp = self.login(session, username, "pass")
        assert resp.status_code == 200
        assert resp.json()['token'] != ""

        resp = self.deregister(session)
        assert resp.status_code == 200

        resp = self.userExists(session, username)
        assert resp.status_code == 200
        assert ('found' not in resp.json() or resp.json()['found'] == False)
