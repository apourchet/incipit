package auth

import (
	"github.com/apourchet/incipit/lib/auth/credentials"
	"github.com/apourchet/incipit/lib/auth/jwt"
)

type mockAuth struct {
	authV1
}

func NewMockAuthClient() AuthClient {
	c := authV1{credentials.NewMockCredentialClient(), jwt.NewDefaultJWTHandler()}
	m := &mockAuth{c}
	return m
}
