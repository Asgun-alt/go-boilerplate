package constant

import (
	"errors"
)

var (
	ErrPasswordNotMatch Error = errors.New("hashedPassword is not the hash of the given password")
	ErrTokenNotFound    Error = errors.New("token not found")
	ErrInvalidToken     Error = errors.New("invalid token")
	ErrOriginNotFound   Error = errors.New("origin not found")
	ErrMethodNotAllowed Error = errors.New("method not allowed")
)
