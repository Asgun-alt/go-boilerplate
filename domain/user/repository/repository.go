package repository

import (
	"go-boilerplate/domain/user/model"
	"go-boilerplate/infrastructure/database"
	"go-boilerplate/pkg/session"
)

type UserRepository interface {
	Create(sess *session.Session, user *model.User) error
	FindByID(sess *session.Session, userID string) (*model.User, error)
	FindByUsername(sess *session.Session, email string) (*model.User, error)
	FindByEmail(sess *session.Session, email string) (*model.User, error)
	Update(sess *session.Session, user *model.User) error
	Delete(sess *session.Session, userID string) error
	FindByEmailAndPassword(sess *session.Session, email, password string) (*model.User, error)
	UpdateIsEmailVerified(sess *session.Session, email string, isVerified bool) error
}

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}
