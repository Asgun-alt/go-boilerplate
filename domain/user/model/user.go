package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Username        string     `json:"username"`
	Password        string     `json:"password"`
	Email           string     `json:"email"`
	UserID          uuid.UUID  `json:"user_id"`
	IsEmailVerified bool       `json:"is_email_verified"`
}

func (u *User) TableName() string {
	return "users"
}
