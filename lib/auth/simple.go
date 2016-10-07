package auth

import (
	"fmt"

	"github.com/apourchet/incipit/lib/auth/credentials"
	"github.com/apourchet/incipit/lib/auth/jwt"
	"github.com/golang/glog"
)

type simpleAuth struct {
	creds      credentials.CredentialsClient
	jwtHandler jwt.JWTHandler
}

func NewK8sAuthClient() (AuthClient, error) {
	glog.Infof("Using K8sAuthClient")
	creds, err := credentials.NewK8sEtcdCredentialsClient()
	if err != nil {
		return nil, err
	}
	return NewSimpleAuthClient(creds, jwt.NewDefaultJWTHandler()), nil
}

func NewSimpleAuthClient(creds credentials.CredentialsClient, jwtHandler jwt.JWTHandler) AuthClient {
	return &simpleAuth{creds, jwtHandler}
}

func (m *simpleAuth) UserExists(key string) (exists bool, err error) {
	return m.creds.UserExists(key)
}

func (m *simpleAuth) Register(key string, pass string) (err error) {
	return m.creds.Register(key, pass)
}

func (m *simpleAuth) Login(key string, pass string) (token string, err error) {
	check, err := m.creds.CheckCredentials(key, pass)
	if err != nil {
		glog.Errorf("Failed to check credentials %v", err)
		return "", fmt.Errorf("Error checking credentials")
	}
	if !check {
		return "", fmt.Errorf("Wrong credentials")
	}
	token, err = m.jwtHandler.CreateToken(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *simpleAuth) Renew(token string) (newToken string, err error) {
	subject, err := m.jwtHandler.ValidateToken(token)
	if err != nil {
		return "", err
	}
	newToken, err = m.jwtHandler.CreateToken(subject)
	if err != nil {
		glog.Errorf("Failed to create new token %v", err)
		return "", err
	}
	return newToken, nil
}

func (m *simpleAuth) Validate(token string) (id string, err error) {
	return m.jwtHandler.ValidateToken(token)
}

func (m *simpleAuth) Logout(token string) (err error) {
	return m.jwtHandler.InvalidateToken(token)
}

func (m *simpleAuth) Deregister(token string) (err error) {
	subject, err := m.jwtHandler.ValidateToken(token)
	if err != nil {
		return err
	}
	err = m.jwtHandler.InvalidateToken(token)
	if err != nil {
		glog.Errorf("Failed to invalidate token %v", err)
		return err
	}
	return m.creds.Deregister(subject)
}
