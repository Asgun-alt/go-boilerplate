package jwt

import "time"

type Config struct {
	SecretKey   string        `json:"secretKey"`
	ExpiredTime time.Duration `json:"jwtExpiredTime"`
}
