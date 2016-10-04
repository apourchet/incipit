package auth

type AuthClient interface {
	UserExists(key string) (exists bool, err error)
	GetUserId(key string) (id string, err error)
	Register(key string, pass string) (err error)

	Login(key string, pass string) (token string, err error)

	Renew(token string) (newToken string, err error)
	Validate(token string) (valid bool, err error)
	Logout(token string) (err error)
	Deregister(token string) (err error)
	CheckClaim(token string, field string, check func(val string) bool) (checked bool, err error)
	AddClaim(token string, field string, val string) (newToken string, err error)
}

type mockAuth struct{}

func NewMockAuth() AuthClient {
	return &mockAuth{}
}

func (m *mockAuth) UserExists(key string) (exists bool, err error) {
	return false, nil
}

func (m *mockAuth) GetUserId(key string) (id string, err error) {
	return "123", nil
}

func (m *mockAuth) Register(key string, pass string) (err error) {
	return nil
}

func (m *mockAuth) Login(key string, pass string) (token string, err error) {
	return "", nil
}

func (m *mockAuth) Renew(token string) (newToken string, err error) {
	return token, nil
}

func (m *mockAuth) Validate(token string) (valid bool, err error) {
	return true, nil
}

func (m *mockAuth) Logout(token string) (err error) {
	return nil
}

func (m *mockAuth) Deregister(token string) (err error) {
	return nil
}

func (m *mockAuth) CheckClaim(token string, field string, check func(val string) bool) (checked bool, err error) {
	return true, nil
}

func (m *mockAuth) AddClaim(token string, field string, val string) (newToken string, err error) {
	return token, nil
}
