package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	h := NewDefaultJWTHandler()
	_, err := h.CreateToken("user")
	assert.Nil(t, err)
}

func TestValidate(t *testing.T) {
	h := NewDefaultJWTHandler()
	token, _ := h.CreateToken("user")
	subject, err := h.ValidateToken(token)
	assert.Nil(t, err)
	assert.Equal(t, subject, "user")
}

func TestInvalid(t *testing.T) {
	h := NewDefaultJWTHandler()
	h.SetLifetime(-10)
	token, err := h.CreateToken("user")
	assert.Nil(t, err)
	_, err = h.ValidateToken(token)
	assert.NotNil(t, err)
}

func TestInvalidate(t *testing.T) {
	h := NewDefaultJWTHandler()
	token, _ := h.CreateToken("user")
	_, err := h.ValidateToken(token)
	assert.Nil(t, err)
	err = h.InvalidateToken(token)
	assert.Nil(t, err)
	_, err = h.ValidateToken(token)
	assert.NotNil(t, err)
}
