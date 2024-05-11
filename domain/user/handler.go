package user

import (
	"go-boilerplate/domain/user/model"
	userusecase "go-boilerplate/domain/user/usecase"
	"go-boilerplate/pkg/response"
	"go-boilerplate/pkg/session"
	"go-boilerplate/pkg/vo"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	CreateUser(e echo.Context) error
}

type handler struct {
	userUcase userusecase.UserUsecase
}

func NewHandler(uc userusecase.UserUsecase) Handler {
	return &handler{
		userUcase: uc,
	}
}

func (h *handler) CreateUser(e echo.Context) error {
	sess := session.GetSession(e)

	request := new(model.CreateUserRequest)
	err := vo.Bind(sess, request)
	if err != nil {
		return response.Error(sess, err)
	}

	err = h.userUcase.CreateUser(sess, request)
	if err != nil {
		return response.Error(sess, err)
	}

	return response.OK(sess, nil)
}
