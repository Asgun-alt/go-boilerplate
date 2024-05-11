package http

import (
	"go-boilerplate/delivery/container"
	auth "go-boilerplate/domain/auth"
)

type Handler struct {
	auth auth.Handler
}

func NewHandler(cont *container.Container) *Handler {
	return &Handler{
		auth: auth.NewHandler(cont.AuthUserUsecase),
	}
}
