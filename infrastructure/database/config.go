package database

import "time"

type Config struct {
	Dialect         string        `json:"dialect"`
	Host            string        `json:"host"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	DBName          string        `json:"dbName"`
	MaxOpenConn     int           `json:"maxOpenConn"`
	MaxIddleConn    int           `json:"maxIddleConn"`
	MaxLifeTimeConn time.Duration `json:"maxLifeTimeConn"`
}
