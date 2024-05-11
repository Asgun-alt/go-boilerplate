package redis

import (
	"encoding/json"
	"fmt"
	"go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"

	"github.com/redis/go-redis/v9"
)

func (r client) Get(sess *session.Session, key string) (val string, err error) {
	val, err = r.rdb.Get(sess.Ctx, key).Result()
	if err != nil {
		if r.Nil(err) {
			return
		}

		err = fmt.Errorf(ErrGetDataRedis.Error(), err)
		sess.SetError(constant.ErrorDatabase, err)
		return
	}

	return
}

func (r client) Set(sess *session.Session, req SetValue) error {
	value, err := json.Marshal(req.Value)
	if err != nil {
		err = fmt.Errorf(ErrSetDataRedis.Error(), err)
		sess.SetError(constant.ErrorDatabase, err)
		return err
	}

	if err := r.rdb.Set(sess.Ctx,
		req.Key,
		value,
		req.TTL,
	).Err(); err != nil {
		if r.Nil(err) {
			return err
		}

		err = fmt.Errorf(ErrSetDataRedis.Error(), err)
		sess.SetError(constant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r client) SetNX(sess *session.Session, req SetValue) error {
	if err := r.rdb.SetNX(sess.Ctx,
		req.Key,
		req.Value,
		req.TTL,
	).Err(); err != nil {
		if r.Nil(err) {
			return err
		}

		err = fmt.Errorf(ErrSetDataRedis.Error(), err)
		sess.SetError(constant.ErrorDatabase, err)
		return err
	}

	return nil
}

func (r client) Del(sess *session.Session, key string) error {
	if err := r.rdb.Del(sess.Ctx,
		key,
	).Err(); err != nil {
		if r.Nil(err) {
			return err
		}

		err = fmt.Errorf(ErrDelDataRedis.Error(), err)
		sess.SetError(constant.ErrorDatabase, err)
		return err
	}

	return nil
}

//nolint:gosimple // keep it like this!
func (r client) Nil(err error) bool {
	if err == redis.Nil {
		return true
	}
	return false
}
