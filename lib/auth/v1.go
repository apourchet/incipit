package auth

import (
	"github.com/apourchet/incipit/lib/auth/credentials"
	"github.com/apourchet/incipit/lib/auth/jwt"
	"github.com/golang/glog"
)

type authV1 struct {
	creds      credentials.CredentialsClient
	jwtHandler jwt.JWTHandler
}

func NewAuthClientV1() (AuthClient, error) {
	glog.Infof("Auth: Using AuthClientV1")
	creds, err := credentials.GetDefaultClient()
	if err != nil {
		glog.Infof("Failed to get EtcdCredentialsClient")
		return nil, err
	}
	a := &authV1{}
	a.creds = creds
	a.jwtHandler = jwt.GetDefaultHandler()
	return a, nil
}

func NewAuthClientV1Fatal() AuthClient {
	c, err := NewAuthClientV1()
	if err != nil {
		glog.Fatalf("Auth: Failed to get AuthClientV1: %v", err)
	}
	return c
}

func (a *authV1) UserExists(key string) (exists bool, err error) {
	glog.Infof("Auth: UserExists: %s", key)
	exists, err = a.creds.UserExists(key)
	return exists, err
}

func (m *authV1) Register(key string, pass string) (ok bool, err error) {
	glog.Infof("Auth: Register: %s, %s", key, pass)
	return m.creds.Register(key, pass)
}

func (m *authV1) Login(key string, pass string) (token string, valid bool, err error) {
	glog.Infof("Auth: Login: %s, %s", key, pass)
	check, err := m.creds.CheckCredentials(key, pass)
	if err != nil {
		glog.Errorf("Auth: Failed to check credentials %v", err)
		return "", false, err
	}
	glog.Infof("Auth: Checked credentials, match: %t", check)
	if !check {
		return "", false, nil
	}
	token, err = m.jwtHandler.CreateToken(key)
	if err != nil {
		return "", false, err
	}
	return token, true, nil
}

func (m *authV1) Validate(token string) (id string, valid bool, err error) {
	return m.jwtHandler.ValidateToken(token)
}

func (m *authV1) Logout(token string) (err error) {
	return m.jwtHandler.InvalidateToken(token)
}

func (m *authV1) Deregister(token string) (err error) {
	subject, valid, err := m.jwtHandler.ValidateToken(token)
	if err != nil || !valid {
		glog.Infof("Auth: Could not validate token: %v %t", err, valid)
		return err
	}
	err = m.jwtHandler.InvalidateToken(token)
	if err != nil {
		glog.Errorf("Auth: Failed to invalidate token %v", err)
		return err
	}

	glog.Infof("Auth: Token was validated %s", subject)
	return m.creds.Deregister(subject)
}
