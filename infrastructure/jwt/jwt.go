package jwt

import "github.com/labstack/echo/v4"

type IJWT interface {
	GenerateToken(customClaim *CustomClaims) (string, int, error)
	VerifyToken(tokenStr string) (*CustomClaims, error)
	RefreshToken(accessTokenString string) (string, int, error)
	VerifyTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type ijwt struct {
	cfgJWT Config
}

func New(
	cJWT Config,
) IJWT {
	return &ijwt{
		cfgJWT: cJWT,
	}
}
