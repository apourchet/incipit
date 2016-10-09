package auth

import (
	"fmt"

	"github.com/apourchet/incipit/lib/auth/credentials"
	"github.com/apourchet/incipit/lib/auth/jwt"
	"github.com/golang/glog"
)

type authV1 struct {
	creds      credentials.CredentialsClient
	jwtHandler jwt.JWTHandler
}

func NewAuthClientV1() (AuthClient, error) {
	glog.Infof("Using AuthClientV1")
	creds, err := credentials.NewEtcdCredentialsClient()
	if err != nil {
		return nil, err
	}
	return &authV1{creds, jwt.NewDefaultJWTHandler()}, nil
}

func NewAuthClientV1Fatal() AuthClient {
	c, err := NewAuthClientV1()
	if err != nil {
		glog.Fatalf("Failed to get AuthClientV1: %v", err)
	}
	return c
}

func (m *authV1) UserExists(key string) (exists bool, err error) {
	return m.creds.UserExists(key)
}

func (m *authV1) Register(key string, pass string) (err error) {
	return m.creds.Register(key, pass)
}

func (m *authV1) Login(key string, pass string) (token string, err error) {
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

func (m *authV1) Renew(token string) (newToken string, err error) {
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

func (m *authV1) Validate(token string) (id string, err error) {
	return m.jwtHandler.ValidateToken(token)
}

func (m *authV1) Logout(token string) (err error) {
	return m.jwtHandler.InvalidateToken(token)
}

func (m *authV1) Deregister(token string) (err error) {
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
