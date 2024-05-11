package model

import (
	userModel "go-boilerplate/domain/user/model"
)

type UserRegistrationRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=12,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=12"`
}

func (m UserRegistrationRequest) ToCreateUserRequest(
	username,
	password string,
) *userModel.CreateUserRequest {
	return &userModel.CreateUserRequest{
		Email:           m.Email,
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
	}
}

type UserRegistrationResponse struct {
	Message string `json:"message"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=12"`
}

func (m UserLoginRequest) ToGetUserRequest(password string) *userModel.GetUserRequest {
	return &userModel.GetUserRequest{
		Email:    m.Email,
		Password: password,
	}
}

type UserLoginResponse struct {
	Token      string `json:"token"`
	Expired    int    `json:"expired"`
	IsVerified bool   `json:"isVerified"`
}

type UserVerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,min=6,max=6"`
}

type UserVerifyResponse struct {
	Mesage  string `json:"message,omitempty"`
	Token   string `json:"token"`
	Expired int    `json:"expired"`
}

type ResendVerificationCode struct {
	Email string `json:"email" validate:"required,email"`
}
