package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"

	"google.golang.org/grpc/metadata"

	"golang.org/x/net/context"
)

const (
	CookieName   = "auth-token"
	CookieMaxAge = 0
)

func GetToken(ctx context.Context) (token string, err error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("Found no token in request")
	}

	authArr, ok := md["authorization"]
	if !ok || len(authArr) != 1 {
		return "", fmt.Errorf("Found no token in request")
	}

	bearerArr := strings.Split(authArr[0], "Bearer ")
	if len(bearerArr) != 2 {
		return "", fmt.Errorf("Malformed bearer header: %s", authArr[0])
	}
	token = bearerArr[1]
	glog.Infof("Got token: %s", token)
	return token, nil
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
