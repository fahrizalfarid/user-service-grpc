package constant

import "errors"

var (
	ErrUsernameNotExists = errors.New("username doesn't exist")
	ErrUsernameExists    = errors.New("username exist")
	ErrEmailNotExists    = errors.New("email doesn't exist")
	ErrEmailExists       = errors.New("email exist")
	ErrSrvNotAvailable   = errors.New("service is not available")
	ErrAuth              = errors.New("token is missing")
)
