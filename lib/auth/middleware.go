package auth

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	// Gets token from request and renews it
	RenewToken() func(*gin.Context)

	// Checks the validity of the token
	// Returns 401 if it is not valid
	CheckToken() func(*gin.Context)

	// Checks the claim on the token
	// Returns 401 if that claim is absent or not satisfied
	CheckClaim(field string, check func(val string) bool) func(*gin.Context)

	// Sets a field on the context to the user id
	// Provided by the AuthClient it uses
	// This is so that endpoints have access to that userid
	SetUserId() func(*gin.Context)
}

type mockAuthMW struct {
	authClient AuthClient
}

func NewMockAuthMW() AuthMiddleware {
	return &mockAuthMW{NewMockAuth()}
}

func (m *mockAuthMW) RenewToken() func(*gin.Context) {
	return func(c *gin.Context) {}
}

func (m *mockAuthMW) CheckToken() func(*gin.Context) {
	return func(c *gin.Context) {}
}

func (m *mockAuthMW) CheckClaim(field string, check func(val string) bool) func(*gin.Context) {
	return func(c *gin.Context) {}
}

func (m *mockAuthMW) SetUserId() func(*gin.Context) {
	return func(c *gin.Context) {}
}
