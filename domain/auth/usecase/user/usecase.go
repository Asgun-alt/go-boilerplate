package authuser

import (
	"go-boilerplate/config"
	"go-boilerplate/domain/auth/model"
	"go-boilerplate/domain/auth/repository"
	"go-boilerplate/infrastructure/jwt"
	"go-boilerplate/pkg/session"

	userusecase "go-boilerplate/domain/user/usecase"
)

type Usecase interface {
	RegisterUser(
		sess *session.Session,
		req *model.UserRegistrationRequest,
	) (response *model.UserRegistrationResponse, err error)
	LoginUser(
		sess *session.Session,
		req *model.UserLoginRequest,
	) (response *model.UserLoginResponse, err error)
	SendVerificationCode(
		sess *session.Session,
		username, email string,
	) error
	VerifyUser(
		sess *session.Session,
		req *model.UserVerifyRequest,
	) (*model.UserVerifyResponse, error)
}

type usecase struct {
	authRepository    repository.AuthRepository
	userUsecase       userusecase.UserUsecase
	cfgAuthCredential config.AuthCredentialConfig
	cfgJwt            jwt.Config
	SmtpConfig        config.SMPTConfig
	jwt               jwt.IJWT
}

func NewAuthUserUsecase(
	authRepo repository.AuthRepository,
	cAuthCred config.AuthCredentialConfig,
	cJwt jwt.Config,
	smptConfig config.SMPTConfig,
	userUc userusecase.UserUsecase,
	ijwt jwt.IJWT,
) Usecase {
	return &usecase{
		authRepository:    authRepo,
		cfgAuthCredential: cAuthCred,
		cfgJwt:            cJwt,
		SmtpConfig:        smptConfig,
		userUsecase:       userUc,
		jwt:               ijwt,
	}
}
