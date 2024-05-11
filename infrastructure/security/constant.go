package security

//nolint:gosec // its just message, not credential
const (
	ErrEncryptPassword = "error to encrypt password, err: %s"
	ErrDecryptPassword = "error to decrypt password, err: %s"

	ErrEncrypt = "error to encrypt, err: %s"
	ErrDecrypt = "error to decrypt, err: %s"
)
