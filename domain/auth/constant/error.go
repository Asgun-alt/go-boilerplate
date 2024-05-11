package constant

import "errors"

type Error error

var (
	UserNotFound       Error = errors.New("user not found")
	WrongEmailPassword Error = errors.New("wrong email or password")

	ErrRequestOTP  = "error email %s the maximum limit for requesting an otp"
	ErrVerifyOTP   = "error email %s the maximum limit for verifying an otp"
	ErrInvalidOTP  = "error invalid otp"
	ErrOTPNotFound = "error otp not found"
)
