package auth

import (
	"fmt"
	"strconv"
)

// Tokens are valid forever
// Any claim is true
// All in memory
// Not thread safe

type mockAuth struct {
	users  map[string]user
	tokens map[string]bool
	currId int
}

func NewMockAuthClient() AuthClient {
	users := make(map[string]user)
	tokens := make(map[string]bool)
	return &mockAuth{users, tokens, 0}
}

func (m *mockAuth) UserExists(key string) (exists bool, err error) {
	_, ok := m.users[key]
	return ok, nil
}

func (m *mockAuth) Register(key string, pass string) (err error) {
	m.users[key] = user{strconv.Itoa(m.currId), key, pass}
	m.currId += 1
	return nil
}

func (m *mockAuth) Login(key string, pass string) (token string, err error) {
	user, ok := m.users[key]
	if !ok {
		return "", fmt.Errorf("User not found")
	}
	if user.password == pass {
		m.tokens[token] = true
		return key, nil
	}
	return "", fmt.Errorf("Wrong credentials")
}

func (m *mockAuth) GetUserId(token string) (id string, err error) {
	valid, err := m.Validate(token)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("Invalid token")
	}
	user, found := m.users[token]
	if !found {
		return "", fmt.Errorf("User not found")
	}
	return user.id, nil
}

func (m *mockAuth) Renew(token string) (newToken string, err error) {
	return token, nil
}

func (m *mockAuth) Validate(token string) (valid bool, err error) {
	_, found1 := m.users[token]
	_, found2 := m.tokens[token]
	if !found1 || !found2 {
		return false, nil
	}
	return true, nil
}

func (m *mockAuth) Logout(token string) (err error) {
	delete(m.tokens, token)
	return nil
}

func (m *mockAuth) Deregister(token string) (err error) {
	delete(m.users, token)
	return
}

func (m *mockAuth) CheckClaim(token string, field string, check func(val string) bool) (checked bool, err error) {
	return true, nil
}

func (m *mockAuth) AddClaim(token string, field string, val string) (newToken string, err error) {
	return token, nil
}
