package userusecase

import (
	"go-boilerplate/domain/user/model"
	"go-boilerplate/domain/user/repository"

	"go-boilerplate/pkg/session"
)

type UserUsecase interface {
	CreateUser(sess *session.Session, req *model.CreateUserRequest) error
	GetUser(sess *session.Session, req *model.GetUserRequest) (*model.User, error)
	GetUserByEmail(sess *session.Session, email string) (*model.User, error)
	GetUserByUserID(sess *session.Session, userID string) (*model.User, error)
	UpdateVerifiedUser(sess *session.Session, email string, isVerified bool) error
}

type userUsecase struct {
	repository repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repo,
	}
}
