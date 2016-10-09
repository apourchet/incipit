package credentials

import "github.com/apourchet/incipit/lib/utils"

type CredentialsClient interface {
	// Whether or not the user already exists
	// Error for internal only
	UserExists(key string) (exists bool, err error)

	// Errors if user exists
	// Errors for internal error
	Register(key string, pass string) (err error)

	// Error for internal only
	// Return "" if user not found
	GetUserId(key string) (id string, err error)

	// Return false if user not found
	// Error for internal only
	CheckCredentials(key string, pass string) (check bool, err error)

	// Error for internal only
	Deregister(key string) (err error)
}

func GetDefaultClient() (CredentialsClient, error) {
	if utils.InKubernetes() {
		return NewEtcdCredentialsClient()
	}
	return NewMockCredentialClient(), nil
}
