package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Aud        string `json:"aud"`
	Name       string `json:"name"`
	Issuer     string `json:"issuer"`
	UserID     string `json:"user_id"`
	IssuerAt   int64  `json:"iat"`
	ExpiredAt  int64  `json:"exp"`
	jwt.Claims `json:"-,omitempty"`
}
