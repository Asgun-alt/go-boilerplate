package redis

import (
	"context"
	"fmt"
	"go-boilerplate/pkg/session"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Get(sess *session.Session, key string) (val string, err error)
	Set(sess *session.Session, req SetValue) error
	SetNX(sess *session.Session, req SetValue) error
	Del(sess *session.Session, key string) error
	Nil(err error) bool
}

type client struct {
	rdb *redis.Client
}

func New(cfg Config) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(fmt.Sprintf("can't connect to db: %s", err.Error()))
	}

	return client{
		rdb: rdb,
	}
}
