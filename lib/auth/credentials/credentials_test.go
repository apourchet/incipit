package credentials

import (
	"testing"

	"github.com/apourchet/incipit/lib/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func credentialClients() []CredentialsClient {
	if utils.InKubernetes() {
		c, _ := NewEtcdCredentialsClient()
		return []CredentialsClient{
			NewMockCredentialClient(),
			c,
		}
	}
	return []CredentialsClient{
		NewMockCredentialClient(),
	}
}

func TestUserExists(t *testing.T) {
	for _, c := range credentialClients() {
		k := uuid.NewV4().String()
		e, _ := c.UserExists(k)
		assert.Equal(t, e, false)

		c.Register(k, "pass")
		e, _ = c.UserExists(k)
		assert.Equal(t, e, true)

		c.Deregister(k)
		e, _ = c.UserExists(k)
		assert.Equal(t, e, false)
	}
}

func TestCheckCredentials(t *testing.T) {
	for _, c := range credentialClients() {
		k := uuid.NewV4().String()
		k1 := uuid.NewV4().String()
		c.Register(k, "pass")
		check, _ := c.CheckCredentials(k, "pass")
		assert.Equal(t, check, true)
		check, _ = c.CheckCredentials(k, "pass1")
		assert.Equal(t, check, false)
		check, _ = c.CheckCredentials(k1, "pass")
		assert.Equal(t, check, false)
		c.Deregister(k)
		check, _ = c.CheckCredentials(k, "pass")
		assert.Equal(t, check, false)
	}
}

// A credential client should return an error on
// user key conflict
func TestDoubleRegister(t *testing.T) {
	for _, c := range credentialClients() {
		k := uuid.NewV4().String()
		c.Register(k, "pass")
		ok, err := c.Register(k, "pass")
		assert.Nil(t, err)
		assert.False(t, ok)
	}
}
