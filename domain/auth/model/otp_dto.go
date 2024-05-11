package model

type OTPRedisRequest struct {
	OTP   string `json:"otp"`
	Email string `json:"email"`
	Count int    `json:"count"`
}

type OTPRedisVerifyRequest struct {
	Email string `json:"email"`
	Count int    `json:"count"`
}

type OTPVerificationEmailPayload struct {
	Username  string
	UserEmail string
	OTPCode   string
}
