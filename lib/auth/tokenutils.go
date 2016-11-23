package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CookieName   = "auth-token"
	CookieMaxAge = 0
)

func SetToken(c *gin.Context, token string) error {
	cookie := newCookieFromToken(token)
	http.SetCookie(c.Writer, cookie)
	return nil
}

func DeleteToken(c *gin.Context) error {
	cookie := newExpiredCookie()
	http.SetCookie(c.Writer, cookie)
	return nil
}

func GetToken(c *gin.Context) (token string, err error) {
	token, err = tokenFromHeader(c)
	if err == nil {
		return token, nil
	}

	// TODO warning here, client should not be using the cookies
	cookie, err := c.Request.Cookie(CookieName)
	if err == nil {
		return cookie.Value, nil
	}
	return "", err
}

func tokenFromHeader(c *gin.Context) (token string, err error) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		return "", nil
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return "", nil
	}
	return strings.Split(header, "Bearer ")[1], nil
}

func newCookieFromToken(token string) *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = CookieName
	cookie.MaxAge = CookieMaxAge
	cookie.Value = token
	return cookie
}

func newExpiredCookie() *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = CookieName
	cookie.MaxAge = -1
	cookie.Value = ""
	return cookie
}
