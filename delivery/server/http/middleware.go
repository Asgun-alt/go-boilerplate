package http

import (
	"bytes"
	"fmt"
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/ratelimiter"
	"go-boilerplate/pkg/response"
	"go-boilerplate/pkg/session"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"go-boilerplate/pkg/utils"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := session.GetSession(c)

		request := c.Request()
		contentType := request.Header.Get("Content-Type")

		var body = make(map[string]interface{}, 0)
		if strings.Contains(contentType, "multipart/form-data") {
			err := request.ParseMultipartForm(32 << 20)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse multipart form: "+err.Error())
			}
			for key, values := range request.MultipartForm.Value {
				if len(values) > 1 {
					body[key] = values
				} else if len(values) == 1 {
					body[key] = values[0]
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "failed reading request body")
			}
			defer request.Body.Close()
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			body, _ = utils.BytesToMap(bodyBytes)
		}

		sess.SetBodyRequest(body)
		sess.SetIP(c.RealIP())

		sess.SetInfo("incoming request")
		c.Set(constant.SessionKey, sess)
		return next(c)
	}
}

func PanicHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := session.GetSession(c)
		defer func() {
			if r := recover(); r != nil {
				var err error
				for i := 1; ; i++ {
					pc, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}

					currentDirectory, _ := os.Getwd()
					if strings.HasPrefix(file, currentDirectory) {
						relativePath := strings.TrimPrefix(file, currentDirectory+"/")
						function := runtime.FuncForPC(pc)
						fnNames := strings.Split(function.Name(), ".")
						funcName := strings.Join(fnNames[len(fnNames)-1:], "")
						panicMessage := fmt.Sprintf("%v", r)

						err = fmt.Errorf("panic in function %s at %s:%d, message: %s",
							funcName,
							relativePath,
							line,
							panicMessage)
						sess.SetError(constant.ErrorPanic, err)
						break
					}
				}

				resp := response.Response{
					Status:  response.ResponseStatusnknown,
					Message: response.ResponseErrorUnknown,
				}

				statusCode := http.StatusInternalServerError
				sess.SetStatusCode(statusCode)
				sess.SetResponse(statusCode, resp, err)
				_ = c.JSON(statusCode, resp)
			}
		}()
		return next(c)
	}
}

func RateLimiter(rl *ratelimiter.RateLimiter, env string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if env != "dev" {
				sess := session.GetSession(c)
				ip := c.RealIP()
				if !rl.Allow(ip) {
					err := fmt.Errorf("to many request from ip %s", ip)
					sess.SetError(constant.ErrorToManyRequest, err)
					return response.Error(sess, err)
				}
			}

			return next(c)
		}
	}
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	sess := session.GetSession(c)
	code := http.StatusInternalServerError
	msg := constant.ErrorInteral

	// Check for a HTTP error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Message != nil {
			msg = constant.Message(he.Message.(string))
		}
	}

	// Customize error message for CORS errors
	if code == http.StatusForbidden && strings.Contains(strings.ToLower(string(msg)), "forbidden") {
		msg = constant.ErrorForbidden
	}

	if code == http.StatusNotFound && strings.Contains(strings.ToLower(string(msg)), "not found") {
		msg = constant.ErrorNotFound
	}

	if code == http.StatusMethodNotAllowed && strings.Contains(strings.ToLower(string(msg)), "method not allowed") {
		msg = constant.ErrorMethodNotAllowed
	}

	if code == http.StatusTooManyRequests && strings.Contains(strings.ToLower(string(msg)), "to many request") {
		msg = constant.ErrorToManyRequest
	}

	response := response.Response{
		Status:  "error",
		Message: response.Message(msg),
		Data:    nil,
	}

	sess.SetResponse(code, response, err)
	if err := c.JSON(code, response); err != nil {
		sess.SetError(msg, err)
	}
}
