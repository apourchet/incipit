package auth

import "github.com/apourchet/incipit/lib/utils"

type AuthClient interface {
	UserExists(key string) (exists bool, err error)
	Register(key string, pass string) (err error)
	Login(key string, pass string) (token string, err error)

	Renew(token string) (newToken string, err error)
	Validate(token string) (id string, err error)

	Logout(token string) (err error)
	Deregister(token string) (err error)
}

func GetDefaultClient() AuthClient {
	if utils.InKubernetes() {
		NewAuthClientV1Fatal()
	}
	return NewMockAuthClient()
}
