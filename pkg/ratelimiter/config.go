package ratelimiter

import "time"

type Config struct {
	Limit int           `json:"limit"`
	TTL   time.Duration `json:"ttl"`
}
