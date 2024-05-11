package authuser

import (
	"go-boilerplate/domain/auth/constant"
	"go-boilerplate/infrastructure/jwt"
	"go-boilerplate/pkg/session"

	pkgconstant "go-boilerplate/pkg/constant"
)

func (u usecase) generateToken(
	sess *session.Session,
	userID,
	username string,
) (string, int, error) {
	tokenModel := &jwt.CustomClaims{
		Aud:    constant.Aud,
		Name:   username,
		Issuer: constant.Issuer,
		UserID: userID,
	}

	accessToken, exp, err := u.jwt.GenerateToken(tokenModel)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return "", 0, err
	}

	err = u.authRepository.SetAccessTokenToRedis(sess,
		accessToken,
		userID,
		u.cfgJwt.ExpiredTime)
	if err != nil {
		return "", 0, err
	}

	return accessToken, exp, err
}
