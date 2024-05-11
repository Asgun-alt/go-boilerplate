package model

type CreateUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=8,max=12"`
	Password        string `json:"password" validate:"required,min=8,max=12,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=12"`
}

type GetUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=12"`
}
