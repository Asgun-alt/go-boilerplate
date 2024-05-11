package userusecase

import (
	"go-boilerplate/domain/user/constant"
	"go-boilerplate/domain/user/model"
	"go-boilerplate/pkg/utils"
	"time"

	pkgconstant "go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *userUsecase) CreateUser(
	sess *session.Session,
	req *model.CreateUserRequest,
) error {
	if err := s.checkUsernameAndEmail(sess,
		req.Username,
		req.Email,
	); err != nil {
		if err == constant.ErrEmailExist || err == constant.ErrUsernameExist {
			err = constant.ErrUsernameOrEmailExist
			sess.SetError(pkgconstant.ErrorDupCheck, err)
			return err
		}

		return err
	}

	hashedPassword, err := utils.GeneratePassword(
		req.ConfirmPassword,
		bcrypt.DefaultCost)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return err
	}

	user := &model.User{
		UserID:    uuid.New(),
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repository.Create(sess, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userUsecase) checkUsernameAndEmail(
	sess *session.Session,
	username,
	email string,
) error {
	existingUser, err := s.repository.FindByUsername(sess, username)
	if err != nil && err != constant.ErrUserNotFound {
		return err
	}

	if existingUser != nil {
		err := constant.ErrUsernameExist
		return err
	}

	existingUser, err = s.repository.FindByEmail(sess, email)
	if err != nil && err != constant.ErrUserNotFound {
		return err
	}

	if existingUser != nil {
		err := constant.ErrEmailExist
		return err
	}

	return nil
}

func (s *userUsecase) GetUser(
	sess *session.Session,
	req *model.GetUserRequest,
) (*model.User, error) {
	hashedPassword, err := utils.GeneratePassword(
		req.Password,
		bcrypt.DefaultCost)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	user, err := s.repository.FindByEmailAndPassword(
		sess,
		req.Email,
		hashedPassword)
	if err != nil {
		return nil, err
	}

	return user,
		nil
}

func (s *userUsecase) GetUserByEmail(
	sess *session.Session,
	email string,
) (*model.User, error) {
	user, err := s.repository.FindByEmail(
		sess,
		email)
	if err != nil && err != constant.ErrUserNotFound {
		return nil, err
	}

	return user,
		nil
}

func (s *userUsecase) GetUserByUserID(
	sess *session.Session,
	userID string,
) (*model.User, error) {
	user, err := s.repository.FindByID(
		sess,
		userID)
	if err != nil && err != constant.ErrUserNotFound {
		return nil, err
	}

	return user,
		nil
}

func (s *userUsecase) UpdateVerifiedUser(
	sess *session.Session,
	email string,
	isVerified bool,
) error {
	err := s.repository.UpdateIsEmailVerified(sess, email, isVerified)
	if err != nil {
		return err
	}

	return nil
}
