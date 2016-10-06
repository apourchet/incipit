package auth

import (
	"fmt"
	"strconv"
)

type CredentialsClient interface {
	UserExists(key string) (exists bool, err error)
	Register(key string, pass string) (err error)
	GetUserId(key string) (id string, err error)
	CheckCredentials(key string, pass string) (check bool, err error)
	Deregister(key string) (err error)
}

type mockCredentialClient struct {
	users  map[string]user
	currId int
}

type user struct {
	id       string
	username string
	password string
}

func NewMockCredentialClient() CredentialsClient {
	m := &mockCredentialClient{}
	m.users = make(map[string]user)
	m.currId = 0
	return m
}

func (m *mockCredentialClient) UserExists(key string) (exists bool, err error) {
	_, ok := m.users[key]
	return ok, nil
}

func (m *mockCredentialClient) Register(key string, pass string) (err error) {
	m.users[key] = user{strconv.Itoa(m.currId), key, pass}
	m.currId += 1
	return nil
}

func (m *mockCredentialClient) GetUserId(key string) (id string, err error) {
	user, ok := m.users[key]
	if !ok {
		return "", fmt.Errorf("User not found")
	}
	return user.id, nil
}

func (m *mockCredentialClient) CheckCredentials(key, pass string) (check bool, err error) {
	user, ok := m.users[key]
	if !ok {
		return false, fmt.Errorf("User not found")
	}
	return user.password == pass, nil
}

func (m *mockCredentialClient) Deregister(key string) error {
	delete(m.users, key)
	return nil
}
