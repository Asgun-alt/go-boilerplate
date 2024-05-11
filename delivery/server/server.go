package server

import (
	"context"
	"fmt"
	"go-boilerplate/delivery/container"
	"go-boilerplate/delivery/server/http"
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"
	"go-boilerplate/pkg/vo"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(cont *container.Container) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.NewSession(c)
			c.Set(constant.SessionKey, sess)

			lang := strings.ToLower(c.Request().Header.Get("Accept-Language"))
			if strings.Contains(lang, "id") {
				lang = "id"
			} else if strings.Contains(lang, "en") {
				lang = "en"
			} else {
				lang = "id"
			}
			sess.SetLanguage(lang)

			merchantID := c.Request().Header.Get("X-MERCHANT-ID")
			sess.SetMerchantID(merchantID)
			return next(c)
		}
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cont.Config.App.Host,
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders: []string{"Content-Type", "X-MERCHANT-ID"},
	}))

	e.Use(http.RateLimiter(cont.HTTPLimiter, cont.Config.App.Env))

	e.HTTPErrorHandler = http.CustomHTTPErrorHandler
	e.Validator = vo.NewCustomValidator()

	h := http.NewHandler(cont)
	http.Router(e, h, cont.JWT)

	go func() {
		err := e.Start(fmt.Sprintf(":%d", cont.Config.App.Port))
		if err != nil {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), cont.Config.Feature.GracefullyShutdown)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
