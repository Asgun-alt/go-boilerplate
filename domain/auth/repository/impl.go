package repository

import (
	"encoding/json"
	"fmt"
	"go-boilerplate/domain/auth/constant"
	"go-boilerplate/domain/auth/model"
	"go-boilerplate/infrastructure/redis"
	pkgconstant "go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"
	"time"
)

const (
	KeyAccessToken = "access-token_%s"
)

func (r authRepository) SetAccessTokenToRedis(
	sess *session.Session,
	accessToken,
	userID string,
	ttl time.Duration,
) error {
	req := redis.SetValue{
		Value: accessToken,
		Key:   fmt.Sprintf(KeyAccessToken, userID),
		TTL:   ttl,
	}

	err := r.redis.SetNX(sess, req)
	if err != nil {
		if r.redis.Nil(err) {
			return nil
		}
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r authRepository) GetAccessTokenFromRedis(
	sess *session.Session,
	userID string,
) (string, error) {
	key := fmt.Sprintf(KeyAccessToken, userID)
	val, err := r.redis.Get(sess, key)
	if err != nil {
		if r.redis.Nil(err) {
			return "", nil
		}
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return "", err
	}

	return val, nil
}

func (r authRepository) GetVerifyOTPFromRedis(
	sess *session.Session,
	email string,
) (*model.OTPRedisVerifyRequest, error) {
	otp, err := r.redis.Get(sess, fmt.Sprintf(constant.DefaultVerifyOTPKey, email))
	if err != nil {
		if r.redis.Nil(err) {
			return nil, nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}

	var result model.OTPRedisVerifyRequest
	err = json.Unmarshal([]byte(otp), &result)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	return &result, nil
}

func (r authRepository) SetMaxVerifyOTPToRedis(
	sess *session.Session,
	value model.OTPRedisVerifyRequest,
	ttl time.Duration,
) error {
	reqSet := redis.SetValue{
		Key:   fmt.Sprintf(constant.DefaultVerifyOTPKey, value.Email),
		Value: value,
		TTL:   ttl,
	}

	if err := r.redis.Set(sess, reqSet); err != nil {
		if r.redis.Nil(err) {
			return nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r authRepository) SetOTPToRedis(
	sess *session.Session,
	value model.OTPRedisRequest,
	ttl time.Duration,
) error {
	reqSet := redis.SetValue{
		Key:   fmt.Sprintf(constant.DefaultOTPKey, value.Email),
		Value: value,
		TTL:   ttl,
	}

	if err := r.redis.Set(sess, reqSet); err != nil {
		if r.redis.Nil(err) {
			return nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r authRepository) GetOTPFromRedis(
	sess *session.Session,
	email string,
) (*model.OTPRedisRequest, error) {
	otp, err := r.redis.Get(sess, fmt.Sprintf(constant.DefaultOTPKey, email))
	if err != nil {
		if r.redis.Nil(err) {
			return nil, nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}

	var result model.OTPRedisRequest
	err = json.Unmarshal([]byte(otp), &result)
	if err != nil {
		sess.SetError(pkgconstant.ErrorGeneral, err)
		return nil, err
	}

	return &result, nil
}

func (r authRepository) DeleleteOTPFromRedis(
	sess *session.Session,
	email string,
) error {
	key := fmt.Sprintf(constant.DefaultOTPKey, email)
	if err := r.redis.Del(sess, key); err != nil {
		if r.redis.Nil(err) {
			return nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r authRepository) DeleteVerifyOTPFromRedis(
	sess *session.Session,
	email string,
) error {
	key := fmt.Sprintf(constant.DefaultVerifyOTPKey, email)
	if err := r.redis.Del(sess, key); err != nil {
		if r.redis.Nil(err) {
			return nil
		}

		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}
