package auth

import "github.com/apourchet/incipit/lib/utils"

type AuthClient interface {
	UserExists(key string) (exists bool, err error)
	Register(key string, pass string) (ok bool, err error)
	Login(key string, pass string) (token string, valid bool, err error)

	Validate(token string) (id string, valid bool, err error)

	Logout(token string) (err error)
	Deregister(token string) (err error)
}

func GetDefaultClient() AuthClient {
	if utils.InKubernetes() {
		return NewAuthClientV1Fatal()
	}
	return NewMockAuthClient()
}
