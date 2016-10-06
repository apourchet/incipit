package auth

import (
	"fmt"

	"github.com/apourchet/incipit/lib/auth/credentials"
	"github.com/apourchet/incipit/lib/auth/jwt"
)

type AuthClient interface {
	UserExists(key string) (exists bool, err error)
	Register(key string, pass string) (err error)
	Login(key string, pass string) (token string, err error)

	Renew(token string) (newToken string, err error)
	Validate(token string) (id string, err error)

	Logout(token string) (err error)
	Deregister(token string) (err error)
}

type mockAuth struct {
	creds      credentials.CredentialsClient
	jwtHandler jwt.JWTHandler
}

func NewMockAuthClient() AuthClient {
	m := &mockAuth{}
	m.creds = credentials.NewMockCredentialClient()
	m.jwtHandler = jwt.NewDefaultJWTHandler()
	return m
}

func (m *mockAuth) UserExists(key string) (exists bool, err error) {
	return m.creds.UserExists(key)
}

func (m *mockAuth) Register(key string, pass string) (err error) {
	return m.creds.Register(key, pass)
}

func (m *mockAuth) Login(key string, pass string) (token string, err error) {
	check, err := m.creds.CheckCredentials(key, pass)
	if err != nil {
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

func (m *mockAuth) Renew(token string) (newToken string, err error) {
	subject, err := m.jwtHandler.ValidateToken(token)
	if err != nil {
		return "", err
	}
	newToken, err = m.jwtHandler.CreateToken(subject)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func (m *mockAuth) Validate(token string) (id string, err error) {
	return m.jwtHandler.ValidateToken(token)
}

func (m *mockAuth) Logout(token string) (err error) {
	return m.jwtHandler.InvalidateToken(token)
}

func (m *mockAuth) Deregister(token string) (err error) {
	subject, err := m.jwtHandler.ValidateToken(token)
	if err != nil {
		return err
	}
	err = m.jwtHandler.InvalidateToken(token)
	if err != nil {
		return err
	}
	return m.creds.Deregister(subject)
}
