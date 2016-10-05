package auth

type AuthClient interface {
	UserExists(key string) (exists bool, err error)
	Register(key string, pass string) (err error)

	Login(key string, pass string) (token string, err error)

	GetUserId(token string) (id string, err error)
	Renew(token string) (newToken string, err error)
	Validate(token string) (valid bool, err error)
	Logout(token string) (err error)
	Deregister(token string) (err error)
	CheckClaim(token string, field string, check func(val string) bool) (checked bool, err error)
	AddClaim(token string, field string, val string) (newToken string, err error)
}
