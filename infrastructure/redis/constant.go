package redis

import (
	"errors"
	"go-boilerplate/pkg/constant"
)

var (
	ErrGetDataRedis constant.Error = errors.New("error get data from redis: %s")
	ErrSetDataRedis constant.Error = errors.New("error set data to redis: %s")
	ErrDelDataRedis constant.Error = errors.New("error delete data from redis: %s")
)
