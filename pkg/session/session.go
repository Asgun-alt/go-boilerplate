package session

import (
	"context"
	"go-boilerplate/pkg/constant"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	SessionKey = "sess_key"
)

type Session struct {
	Ctx               context.Context
	EchoCtx           echo.Context
	Error             error
	Response          interface{}
	Request           map[string]interface{}
	Rest              map[string]interface{}
	Message           constant.Message
	IP                string
	Language          string
	DictionaryMessage string
	AccessToken       string
	UserID            string
	MerchantID        string
	Email             string
	StatusCode        int
}

func NewSession(e echo.Context,
) *Session {
	return &Session{
		Ctx:     e.Request().Context(),
		EchoCtx: e,
	}
}

func GetEventSession() *Session {
	return &Session{
		Ctx: context.Background(),
	}
}

func (sess *Session) SetError(message constant.Message, err error) {
	sess.Message = message
	sess.Error = err
}

func (sess *Session) SetInfo(message constant.Message) {
	sess.Message = message
}

func (sess *Session) SetResponse(statusCode int, response interface{}, err error) {
	sess.Response = response
	sess.StatusCode = statusCode

	if err != nil {
		sess.SetError(sess.GetMessage(), err)
	} else {
		sess.SetInfo(sess.GetMessage())
	}
}

func (sess *Session) SetIP(ip string) {
	sess.IP = ip
}

func (sess *Session) SetStatusCode(statusCode int) {
	sess.StatusCode = statusCode
}

func (sess *Session) SetBodyRequest(request map[string]interface{}) {
	sess.Request = request
}

func (sess *Session) SetLanguage(lang string) {
	sess.Language = lang
}

func (sess *Session) SetAccessToken(token string) {
	bearer := "bearer "
	sess.AccessToken = token
	if !strings.Contains(strings.ToLower(token), bearer) {
		sess.AccessToken = bearer + token
	}
}

func (sess *Session) SetUserID(userID string) {
	sess.UserID = userID
}

func (sess *Session) SetMerchantID(merchantID string) {
	sess.MerchantID = merchantID
}

func (sess *Session) SetEmail(email string) {
	sess.Email = email
}

func (sess *Session) GetEventSession() *Session {
	return &Session{
		Ctx: context.Background(),
	}
}

func GetSession(e echo.Context) *Session {
	sess, _ := e.Get(constant.SessionKey).(*Session)
	return sess
}

func (sess *Session) GetMessage() constant.Message {
	return sess.Message
}

func (sess *Session) GetEchoContext() echo.Context {
	return sess.EchoCtx
}

func (sess *Session) GetContext() context.Context {
	return sess.Ctx
}

func (sess *Session) GetIP() string {
	return sess.IP
}

func (sess *Session) GetRest() map[string]interface{} {
	return sess.Rest
}

func (sess *Session) GetError() error {
	return sess.Error
}

func (sess *Session) GetBodyRequest() map[string]interface{} {
	return sess.Request
}

func (sess *Session) GetResponse() interface{} {
	return sess.Response
}

func (sess *Session) GetStatusCode() int {
	return sess.StatusCode
}

func (sess *Session) GetAccessToken() string {
	return sess.AccessToken
}

func (sess *Session) GetUserID() string {
	return sess.UserID
}

func (sess *Session) GetEmail() string {
	return sess.Email
}
