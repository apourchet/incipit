package auth

import (
	"testing"

	"github.com/apourchet/incipit/lib/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func authClients() []AuthClient {
	if utils.InKubernetes() {
		return []AuthClient{
			NewMockAuthClient(),
			NewAuthClientV1Fatal(),
		}
	}
	return []AuthClient{
		NewMockAuthClient(),
	}
}

func TestLogin(t *testing.T) {
	for _, m := range authClients() {
		k := uuid.NewV4().String()
		_, valid, err := m.Login(k, "pass")
		assert.Nil(t, err)
		assert.False(t, valid)
		err = m.Register(k, "pass")
		assert.Nil(t, err)
		token, valid, err := m.Login(k, "pass")
		assert.Nil(t, err)
		assert.True(t, valid)
		assert.NotEqual(t, token, "")
	}
}

func TestDeregister(t *testing.T) {
	for _, m := range authClients() {
		k := uuid.NewV4().String()
		m.Register(k, "pass")
		token, valid, _ := m.Login(k, "pass")
		assert.True(t, valid)
		err := m.Deregister(token + "bad")
		assert.Nil(t, err)
		err = m.Deregister(token)
		assert.Nil(t, err)
		_, valid, err = m.Validate(token)
		assert.Nil(t, err)
		assert.False(t, valid)
	}
}

func TestLogout(t *testing.T) {
	for _, m := range authClients() {
		k := uuid.NewV4().String()

		m.Register(k, "pass")
		token, valid, _ := m.Login(k, "pass")
		assert.True(t, valid)
		err := m.Logout(token + "bad")
		assert.NotNil(t, err)
		err = m.Logout(token)
		assert.Nil(t, err)
		_, valid, err = m.Validate(token)
		assert.Nil(t, err)
		assert.False(t, valid)
	}
}
