package jwt

const (
	ErrGenerateAccessToken     = "error generate access token, err: %s"
	ErrVerifyAccessToken       = "error verifying access token, err: %s"
	ErrUnexpectedSigningMethod = "Unexpected signing method: %v"
	ErrInvalidToken            = "error invalid tokens"
)
