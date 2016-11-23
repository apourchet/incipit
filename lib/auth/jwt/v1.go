package jwt

import (
	"fmt"
	"time"

	"github.com/apourchet/incipit/lib/redis"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type redisJWTHandlerV1 struct {
	store    redis.RedisClient
	lifetime time.Duration
	issuer   string
	key      string
}

func NewRedisJWTHandlerV1() JWTHandler {
	h := &redisJWTHandlerV1{}
	h.store = redis.GetDefaultClient()
	h.lifetime = 10 * time.Minute
	h.issuer = "incipit|redisJWTHandler"
	h.key = "JIUzI1NiIsInR5cCI6IkpXVCJ"
	return h
}

func (h *redisJWTHandlerV1) CreateToken(subject string) (token string, err error) {
	claims := h.createClaims(subject)
	err = h.store.SetExpire(claims.Id, "true", h.lifetime)
	if err != nil {
		return "", err
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString([]byte(h.key))
	return ss, err
}

func (h *redisJWTHandlerV1) ValidateToken(token string) (subject string, valid bool, err error) {
	tok, err := jwt.ParseWithClaims(token, &JWTClaims{}, h.getKeyFromAlg)
	if !tok.Valid {
		return "", false, nil
	}
	claims, ok := tok.Claims.(*JWTClaims)
	if !ok {
		return "", false, nil
	}
	_, found, err := h.store.Get(claims.Id)
	if err != nil || !found {
		return "", false, err
	}
	return claims.Subject, true, nil
}

func (h *redisJWTHandlerV1) InvalidateToken(token string) error {
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
	err = h.store.Delete(claims.Id)
	return nil
}

func (h *redisJWTHandlerV1) SetLifetime(lifetime time.Duration) {
	h.lifetime = lifetime
}

func (h *redisJWTHandlerV1) getKeyFromAlg(token *jwt.Token) (interface{}, error) {
	return []byte(h.key), nil
}

func (h *redisJWTHandlerV1) createClaims(subject string) *JWTClaims {
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
