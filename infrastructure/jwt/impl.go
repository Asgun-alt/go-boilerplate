package jwt

import (
	"fmt"
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/response"
	"go-boilerplate/pkg/session"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (j ijwt) GenerateToken(customClaim *CustomClaims) (token string, expiredAt int, err error) {
	customClaim.IssuerAt = time.Now().Unix()
	customClaim.ExpiredAt = time.Now().Add(j.cfgJWT.ExpiredTime).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)
	token, err = accessToken.SignedString([]byte(j.cfgJWT.SecretKey))
	if err != nil {
		err = fmt.Errorf(ErrGenerateAccessToken, err.Error())
		return "", 0, err
	}

	return token, int(customClaim.ExpiredAt), nil
}

func (i ijwt) VerifyToken(accessTokenString string) (*CustomClaims, error) {
	token, err := jwt.Parse(accessTokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(ErrUnexpectedSigningMethod, token.Header["alg"])
			}

			return []byte(i.cfgJWT.SecretKey), nil
		})

	if err != nil {
		return nil, fmt.Errorf(ErrVerifyAccessToken, err)
	}

	var customClaims CustomClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		customClaims = CustomClaims{
			Aud:       claims["aud"].(string),
			Name:      claims["name"].(string),
			Issuer:    claims["issuer"].(string),
			UserID:    claims["user_id"].(string),
			IssuerAt:  int64(claims["exp"].(float64)),
			ExpiredAt: int64(claims["iat"].(float64)),
		}
		return &customClaims, nil
	}

	return nil, fmt.Errorf(ErrInvalidToken)
}

func (i ijwt) RefreshToken(accessTokenString string) (accessToken string, exp int, err error) {
	claims, err := i.VerifyToken(accessTokenString)
	if err != nil {
		return "", 0, err
	}

	accessToken, exp, err = i.GenerateToken(claims)
	if err != nil {
		return "", 0, err
	}

	return accessToken, exp, nil
}

func (j ijwt) VerifyTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := session.GetSession(c)

		accessToken := c.Request().Header.Get("authorization")
		if accessToken == "" ||
			!strings.Contains(strings.ToLower(accessToken), "bearer ") {
			sess.SetError(constant.ErrorUnauthorized, constant.ErrTokenNotFound)
			return response.Error(sess, constant.ErrTokenNotFound)
		}

		token := strings.Split(accessToken, "Bearer ")
		if len(token) > 1 {
			accessToken = token[1]
			customClaims, err := j.VerifyToken(accessToken)
			if err != nil {
				sess.SetError(constant.ErrorUnauthorized, err)
				return response.Error(sess, err)
			}

			merchantID := c.Request().Header.Get("X-MERCHANT-ID")
			sess.SetMerchantID(merchantID)
			sess.SetUserID(customClaims.UserID)
			sess.SetAccessToken(accessToken)
			c.Set(constant.SessionKey, sess)
			return next(c)
		}

		sess.SetError(constant.ErrorUnauthorized, constant.ErrTokenNotFound)
		return response.Error(sess, constant.ErrTokenNotFound)
	}
}
