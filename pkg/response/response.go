package response

import (
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func BadRequest(sess *session.Session, e echo.Context, err error) error {
	statusCode := http.StatusBadRequest

	msg := sess.Message
	if err != nil {
		msg = constant.Message(err.Error())
	}

	response := Response{
		Status:  ResponseStatusError,
		Message: Message(msg),
		Code:    statusCode,
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		validationErrors := make([]ValidationError, 0)
		for _, err := range err.(validator.ValidationErrors) {
			validationError := ValidationError{
				Field:    strings.ToLower(err.Field()),
				Tag:      strings.ToLower(err.Tag()),
				TagValue: strings.ToLower(err.Param()),
				Value:    err.Value().(string),
			}
			validationErrors = append(validationErrors, validationError)
		}

		response = Response{
			Status:  ResponseStatusError,
			Message: ResponseErrorBadRequest,
			Data:    validationErrors,
			Code:    statusCode,
		}
	}

	go sess.SetResponse(statusCode, response, err)
	return e.JSON(statusCode, response)
}

func OK(sess *session.Session, data interface{}) error {
	e := sess.EchoCtx

	statusCode := http.StatusOK

	msg := ResponseMessageSuccess
	response := Response{
		Status:  ResponseStatusSuccess,
		Message: msg,
		Data:    data,
		Code:    statusCode,
	}

	go sess.SetResponse(statusCode, response, nil)
	return e.JSON(statusCode, response)
}

func Error(sess *session.Session, err error) error {
	e := sess.EchoCtx

	var (
		statusCode int
		status     = ResponseStatusError
	)

	switch {
	case sess.Message == constant.ErrorRequest:
		return BadRequest(sess, e, err)
	case sess.Message == constant.ErrorNotFound:
		statusCode = http.StatusNotFound
	case sess.Message == constant.ErrorUnauthorized:
		statusCode = http.StatusUnauthorized
	case sess.Message == constant.ErrorDupCheck:
		statusCode = http.StatusConflict
	case sess.Message == constant.ErrorForbidden:
		statusCode = http.StatusForbidden
	case sess.Message == constant.ErrorToManyRequest:
		statusCode = http.StatusTooManyRequests
	case sess.Message == constant.ErrorMethodNotAllowed:
		statusCode = http.StatusMethodNotAllowed
	default:
		statusCode = http.StatusInternalServerError
		status = ResponseStatusError
	}

	msg := constant.Message(sess.GetError().Error())
	if sess.GetError().Error() == "" {
		msg = sess.Message
	}

	response := Response{
		Status:  status,
		Message: Message(msg),
		Data:    nil,
		Code:    statusCode,
	}

	go sess.SetResponse(statusCode, response, err)
	return e.JSON(statusCode, response)
}
