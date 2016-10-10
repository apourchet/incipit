package jwt

import (
	"testing"
	"time"

	"github.com/apourchet/incipit/lib/utils"
	"github.com/stretchr/testify/assert"
)

func jwtHandlers() []JWTHandler {
	a := []JWTHandler{
		NewMockJWTHandler(),
	}
	if utils.InKubernetes() {
		a = append(a, NewRedisJWTHandlerV1())
	}
	return a
}

func TestCreate(t *testing.T) {
	for _, h := range jwtHandlers() {
		_, err := h.CreateToken("user")
		assert.Nil(t, err)
	}
}

func TestValidate(t *testing.T) {
	for _, h := range jwtHandlers() {
		token, _ := h.CreateToken("user")
		subject, valid, err := h.ValidateToken(token)
		assert.Nil(t, err)
		assert.Equal(t, subject, "user")
		assert.True(t, valid)
	}
}

func TestInvalid(t *testing.T) {
	for _, h := range jwtHandlers() {
		h.SetLifetime(-10 * time.Second)
		token, err := h.CreateToken("user")
		assert.Nil(t, err)
		_, valid, err := h.ValidateToken(token)
		assert.Nil(t, err)
		assert.False(t, valid)
	}
}

func TestInvalidate(t *testing.T) {
	for _, h := range jwtHandlers() {
		token, _ := h.CreateToken("user")
		_, valid, err := h.ValidateToken(token)
		assert.Nil(t, err)
		assert.True(t, valid)
		err = h.InvalidateToken(token)
		assert.Nil(t, err)
		_, valid, err = h.ValidateToken(token)
		assert.Nil(t, err)
		assert.False(t, valid)
	}
}
