#! /usr/bin/python

import requests

class Auth():
    def __init__(self, ing_addr):
        self.ing_addr = ing_addr

    def newsession(self):
        session = requests.Session()
        session.verify = False
        session.headers.update({
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        })
        return session

    def login(self, username, password):
        session = self.newsession()
        url = 'https://%s/api/v1/auth/login' % self.ing_addr
        url += '?key=%s&pass=%s' % (username, password)

        resp = session.get(url)
        if resp.status_code != 200:
            raise Exception("Login failed: Error: %s" % resp.json()["error"])

        session.headers.update({ 'Authorization': 'Bearer %s' % resp.json()["token"] })
        return session

class Test():
    def __init__(self, ing_addr):
        self.auth = Auth(ing_addr)

    def test(self):
        print "tested at", self.auth.ing_addr


if __name__ == "__main__":
    test = Test('incipit.machine:30002')
    test.test()
