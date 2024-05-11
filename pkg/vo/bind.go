package vo

import (
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"
)

func Bind(sess *session.Session, i interface{}) error {
	e := sess.EchoCtx
	err := e.Bind(i)
	if err != nil {
		sess.SetError(constant.ErrorRequest, err)
		return err
	}

	err = e.Validate(i)
	if err != nil {
		sess.SetError(constant.ErrorRequest, err)
		return err
	}

	return nil
}
