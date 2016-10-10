package jwt

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type mockJWTHandler struct {
	invalids map[string]bool
	lifetime time.Duration
	issuer   string
	key      string
}

func NewMockJWTHandler() JWTHandler {
	h := &mockJWTHandler{}
	h.invalids = make(map[string]bool)
	h.lifetime = 10 * time.Minute
	h.issuer = "incipit|mockJWTHandler"
	h.key = "JIUzI1NiIsInR5cCI6IkpXVCJ"
	return h
}

func (h *mockJWTHandler) CreateToken(subject string) (token string, err error) {
	claims := h.createClaims(subject)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString([]byte(h.key))
	return ss, err
}

func (h *mockJWTHandler) ValidateToken(token string) (subject string, valid bool, err error) {
	tok, err := jwt.ParseWithClaims(token, &JWTClaims{}, h.getKeyFromAlg)
	if !tok.Valid {
		return "", false, nil
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok {
		return "", false, nil
	}
	if _, ok := h.invalids[claims.Id]; ok {
		return "", false, nil
	}
	return claims.Subject, true, nil
}

func (h *mockJWTHandler) InvalidateToken(token string) error {
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

func (h *mockJWTHandler) SetLifetime(lifetime time.Duration) {
	h.lifetime = lifetime
}

func (h *mockJWTHandler) getKeyFromAlg(token *jwt.Token) (interface{}, error) {
	return []byte(h.key), nil
}

func (h *mockJWTHandler) createClaims(subject string) *JWTClaims {
	claims := &JWTClaims{}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Unix() + int64(h.lifetime.Seconds())
	claims.Issuer = h.issuer
	claims.NotBefore = 0
	claims.Audience = "['" + subject + "']"
	claims.Subject = subject
	claims.Id = uuid.NewV4().String()
	return claims
}
