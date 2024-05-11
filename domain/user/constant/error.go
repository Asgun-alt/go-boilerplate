package constant

import (
	"errors"
)

type Error error

var (
	ErrUsernameOrEmailExist Error = errors.New("username or email already exists")
	ErrUsernameExist        Error = errors.New("username already exists")
	ErrEmailExist           Error = errors.New("email already exists")
	ErrUserNotFound         Error = errors.New("user not found")
)
