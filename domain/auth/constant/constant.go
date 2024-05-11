package constant

const (
	SuccessRegisterUser = "register user success"
	UserVerified        = "user has been verified"

	UsernameIncompletedProfileUser     = "user%d"
	UsernameIncompletedProfileCustomer = "user%d"
)

const (
	MagicNumberOTP      = 6
	OTPExpired          = 5
	OTPMaxRequest       = 3
	DefaultOTPKey       = "otp_%s"
	DefaultVerifyOTPKey = "verify_otp_%s"

	Aud    = "Boilerplate-Accessor"
	Issuer = "Boilerplalte-Authority"
)
