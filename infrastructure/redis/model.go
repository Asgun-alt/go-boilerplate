package redis

import "time"

type Config struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	DB       int    `json:"db"`
}

type SetValue struct {
	Value interface{}
	Key   string
	TTL   time.Duration
}
