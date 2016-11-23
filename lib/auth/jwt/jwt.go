package jwt

import (
	"time"

	"github.com/apourchet/incipit/lib/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTHandler interface {
	CreateToken(subject string) (token string, err error)
	ValidateToken(token string) (subject string, validated bool, err error)
	InvalidateToken(token string) (err error)
	SetLifetime(lifetime time.Duration)
}

type JWTClaims struct {
	jwt.StandardClaims
}

func GetDefaultHandler() JWTHandler {
	if utils.InKubernetes() {
		return NewRedisJWTHandlerV1()
	}
	return NewMockJWTHandler()
}
