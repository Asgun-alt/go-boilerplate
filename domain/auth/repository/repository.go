package repository

import (
	"go-boilerplate/domain/auth/model"
	"go-boilerplate/infrastructure/redis"
	"go-boilerplate/pkg/session"
	"time"
)

type AuthRepository interface {
	SetAccessTokenToRedis(
		sess *session.Session,
		accessToken,
		userID string,
		ttl time.Duration,
	) error
	GetAccessTokenFromRedis(
		sess *session.Session,
		userID string,
	) (string, error)
	GetVerifyOTPFromRedis(
		sess *session.Session,
		email string,
	) (*model.OTPRedisVerifyRequest, error)
	SetMaxVerifyOTPToRedis(
		sess *session.Session,
		value model.OTPRedisVerifyRequest,
		ttl time.Duration,
	) error
	SetOTPToRedis(
		sess *session.Session,
		value model.OTPRedisRequest,
		ttl time.Duration,
	) error
	GetOTPFromRedis(
		sess *session.Session,
		email string,
	) (*model.OTPRedisRequest, error)
	DeleleteOTPFromRedis(
		sess *session.Session,
		email string,
	) error
	DeleteVerifyOTPFromRedis(
		sess *session.Session,
		email string,
	) error
}

type authRepository struct {
	redis redis.Client
}

func NewAuthRepository(rdb redis.Client) AuthRepository {
	return &authRepository{
		redis: rdb,
	}
}
