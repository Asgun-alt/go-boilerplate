package http

import (
	"go-boilerplate/infrastructure/jwt"
	netHTTP "net/http"

	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo, h *Handler, j jwt.IJWT) {
	e.Use(Logger)
	e.Use(PanicHandler)

	e.GET("/health-check", func(c echo.Context) error {
		return c.String(netHTTP.StatusOK, "healthy")
	})

	v1 := e.Group("/v1")

	auth := v1.Group("/auth")
	{
		user := auth.Group("/user")
		{
			user.POST("/register", h.auth.RegisterUser)
			user.POST("/login", h.auth.LoginUser)
			user.POST("/verify", h.auth.SendVerificationCode, j.VerifyTokenMiddleware)
			user.POST("/resend-verification", h.auth.ResendVerificationCode, j.VerifyTokenMiddleware)
		}
	}
}
