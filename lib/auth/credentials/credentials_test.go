package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func credentialClients() []CredentialsClient {
	return []CredentialsClient{
		NewMockCredentialClient(),
	}
}

func TestUserExists(t *testing.T) {
	for _, c := range credentialClients() {
		e, _ := c.UserExists("user")
		assert.Equal(t, e, false)

		c.Register("user", "pass")
		e, _ = c.UserExists("user")
		assert.Equal(t, e, true)

		c.Deregister("user")
		e, _ = c.UserExists("user")
		assert.Equal(t, e, false)
	}
}

func TestCheckCredentials(t *testing.T) {
	for _, c := range credentialClients() {
		c.Register("user", "pass")
		check, _ := c.CheckCredentials("user", "pass")
		assert.Equal(t, check, true)
		check, _ = c.CheckCredentials("user", "pass1")
		assert.Equal(t, check, false)
		check, _ = c.CheckCredentials("user1", "pass")
		assert.Equal(t, check, false)
		c.Deregister("user")
		check, _ = c.CheckCredentials("user", "pass")
		assert.Equal(t, check, false)
	}
}

// A credential client should return an error on
// user key conflict
func TestDoubleRegister(t *testing.T) {
	for _, c := range credentialClients() {
		c.Register("user", "pass")
		err := c.Register("user", "pass")
		assert.NotNil(t, err)
	}
}
