package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type JWTHandler interface {
	CreateToken(subject string) (id string, token string, err error)
	ValidateToken(token string) (id string, subject string, err error)
}

type JWTClaims struct {
	jwt.StandardClaims
}

type defaultJWTHandler struct{}

const (
	JWTLifetime   = 10
	JWTIssuer     = "incipit|defaultJWTHandler"
	JWTSigningKey = "JIUzI1NiIsInR5cCI6IkpXVCJ"
)

func NewDefaultJWTHandler() JWTHandler {
	h := &defaultJWTHandler{}
	return h
}

func (h *defaultJWTHandler) CreateToken(subject string) (id string, token string, err error) {
	claims := createClaims(subject)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString([]byte(JWTSigningKey))
	return claims.Id, ss, err
}

func (h *defaultJWTHandler) ValidateToken(token string) (id string, subject string, err error) {
	tok, err := jwt.ParseWithClaims(token, &JWTClaims{}, getKeyFromAlg)
	if !tok.Valid {
		return "", "", fmt.Errorf("Token expired")
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok {
		return "", "", fmt.Errorf("Failed to parse token")
	}
	return claims.Id, claims.Subject, nil
}

func getKeyFromAlg(token *jwt.Token) (interface{}, error) {
	return []byte(JWTSigningKey), nil
}

func createClaims(subject string) *JWTClaims {
	claims := &JWTClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Unix() + JWTLifetime
	claims.Issuer = JWTIssuer
	claims.NotBefore = 0
	claims.Audience = "['" + subject + "']"
	claims.Subject = subject
	claims.Id = uuid.NewV4().String()
	return claims
}
