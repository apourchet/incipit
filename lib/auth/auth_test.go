package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	m := NewMockAuthClient()
	_, err := m.Login("user", "pass")
	assert.NotNil(t, err)
	err = m.Register("user", "pass")
	assert.Nil(t, err)
	token, err := m.Login("user", "pass")
	assert.Nil(t, err)
	assert.NotEqual(t, token, "")
}

func TestRenew(t *testing.T) {
	m := NewMockAuthClient()
	m.Register("user", "pass")
	token, _ := m.Login("user", "pass")
	newToken, err := m.Renew(token)
	assert.Nil(t, err)
	assert.NotEqual(t, newToken, "")
}

func TestDeregister(t *testing.T) {
	m := NewMockAuthClient()
	m.Register("user", "pass")
	token, _ := m.Login("user", "pass")
	err := m.Deregister(token + "bad")
	assert.NotNil(t, err)
	err = m.Deregister(token)
	assert.Nil(t, err)
}

func TestLogout(t *testing.T) {
	m := NewMockAuthClient()
	m.Register("user", "pass")
	token, _ := m.Login("user", "pass")
	err := m.Logout(token + "bad")
	assert.NotNil(t, err)
	err = m.Logout(token)
	assert.Nil(t, err)
	_, err = m.Validate(token)
	assert.NotNil(t, err)
}
