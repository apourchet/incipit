package jwt

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type JWTHandler interface {
	CreateToken(subject string) (token string, err error)
	ValidateToken(token string) (subject string, err error)
	InvalidateToken(token string) (err error)
	SetLifetime(lifetime int64)
}

type JWTClaims struct {
	jwt.StandardClaims
}

type defaultJWTHandler struct {
	invalids map[string]bool
	lifetime int64
	issuer   string
	key      string
}

func NewDefaultJWTHandler() JWTHandler {
	h := &defaultJWTHandler{}
	h.invalids = make(map[string]bool)
	h.lifetime = int64(10)
	h.issuer = "incipit|defaultJWTHandler"
	h.key = "JIUzI1NiIsInR5cCI6IkpXVCJ"
	return h
}

func (h *defaultJWTHandler) CreateToken(subject string) (token string, err error) {
	claims := h.createClaims(subject)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString([]byte(h.key))
	return ss, err
}

func (h *defaultJWTHandler) ValidateToken(token string) (subject string, err error) {
	tok, err := jwt.ParseWithClaims(token, &JWTClaims{}, h.getKeyFromAlg)
	if !tok.Valid {
		return "", fmt.Errorf("Token expired")
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok {
		return "", fmt.Errorf("Failed to parse token")
	}
	if _, ok := h.invalids[claims.Id]; ok {
		return "", fmt.Errorf("Token invalidated")
	}
	return claims.Subject, nil
}

func (h *defaultJWTHandler) InvalidateToken(token string) error {
	tok, err := jwt.ParseWithClaims(token, &JWTClaims{}, h.getKeyFromAlg)
	if err != nil {
		return err
	}
	if !tok.Valid {
		return fmt.Errorf("Token expired")
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok {
		return fmt.Errorf("Failed to parse token")
	}
	h.invalids[claims.Id] = true
	return nil
}

func (h *defaultJWTHandler) SetLifetime(lifetime int64) {
	h.lifetime = lifetime
}

func (h *defaultJWTHandler) getKeyFromAlg(token *jwt.Token) (interface{}, error) {
	return []byte(h.key), nil
}

func (h *defaultJWTHandler) createClaims(subject string) *JWTClaims {
	claims := &JWTClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Unix() + h.lifetime
	claims.Issuer = h.issuer
	claims.NotBefore = 0
	claims.Audience = "['" + subject + "']"
	claims.Subject = subject
	claims.Id = uuid.NewV4().String()
	return claims
}
