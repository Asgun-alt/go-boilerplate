package auth

import (
	"go-boilerplate/domain/auth/model"
	authuser "go-boilerplate/domain/auth/usecase/user"
	pkgconstant "go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/response"
	"go-boilerplate/pkg/session"
	"go-boilerplate/pkg/vo"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	RegisterUser(e echo.Context) error
	LoginUser(e echo.Context) error
	SendVerificationCode(e echo.Context) error
	ResendVerificationCode(e echo.Context) error
}
type handler struct {
	authUserUsecase authuser.Usecase
}

func NewHandler(
	authUserUc authuser.Usecase,
) Handler {
	return &handler{
		authUserUsecase: authUserUc,
	}
}

func (h handler) RegisterUser(e echo.Context) error {
	sess := session.GetSession(e)

	request := new(model.UserRegistrationRequest)
	err := vo.Bind(sess, request)
	if err != nil {
		sess.SetError(pkgconstant.ErrorRequest, err)
		return response.Error(sess, err)
	}

	resp, err := h.authUserUsecase.RegisterUser(sess, request)
	if err != nil {
		return response.Error(sess, err)
	}

	return response.OK(sess, resp)
}

func (h handler) LoginUser(e echo.Context) error {
	sess := session.GetSession(e)

	request := new(model.UserLoginRequest)
	err := vo.Bind(sess, request)
	if err != nil {
		sess.SetError(pkgconstant.ErrorRequest, err)
		return response.Error(sess, err)
	}

	resp, err := h.authUserUsecase.LoginUser(sess, request)
	if err != nil {
		return response.Error(sess, err)
	}

	return response.OK(sess, resp)
}

func (h handler) SendVerificationCode(e echo.Context) error {
	sess := session.GetSession(e)

	request := new(model.UserVerifyRequest)
	err := vo.Bind(sess, request)
	if err != nil {
		sess.SetError(pkgconstant.ErrorRequest, err)
		return response.Error(sess, err)
	}

	resp, err := h.authUserUsecase.VerifyUser(sess, request)
	if err != nil {
		return response.Error(sess, err)
	}

	return response.OK(sess, resp)
}

func (h handler) ResendVerificationCode(e echo.Context) error {
	sess := session.GetSession(e)

	request := new(model.ResendVerificationCode)
	err := vo.Bind(sess, request)
	if err != nil {
		sess.SetError(pkgconstant.ErrorRequest, err)
		return response.Error(sess, err)
	}

	username := ""
	err = h.authUserUsecase.SendVerificationCode(sess, username, request.Email)
	if err != nil {
		return response.Error(sess, err)
	}

	return response.OK(sess, nil)
}
