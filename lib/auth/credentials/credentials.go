package credentials

import (
	"fmt"
	"strconv"
)

type CredentialsClient interface {
	// Whether or not the user already exists
	// Error for internal only
	UserExists(key string) (exists bool, err error)

	// Errors if user exists
	// Errors for internal error
	Register(key string, pass string) (err error)

	// Error for internal only
	// Return "" if user not found
	GetUserId(key string) (id string, err error)

	// Return false if user not found
	// Error for internal only
	CheckCredentials(key string, pass string) (check bool, err error)

	// Error for internal only
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
	_, ok := m.users[key]
	if ok {
		return fmt.Errorf("User already exists")
	}
	m.users[key] = user{strconv.Itoa(m.currId), key, pass}
	m.currId += 1
	return nil
}

func (m *mockCredentialClient) GetUserId(key string) (id string, err error) {
	user, ok := m.users[key]
	if !ok {
		return "", nil
	}
	return user.id, nil
}

func (m *mockCredentialClient) CheckCredentials(key, pass string) (check bool, err error) {
	user, ok := m.users[key]
	if !ok {
		return false, nil
	}
	return user.password == pass, nil
}

func (m *mockCredentialClient) Deregister(key string) error {
	delete(m.users, key)
	return nil
}
