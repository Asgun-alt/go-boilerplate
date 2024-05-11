package ratelimiter

import (
	"sync"
	"time"
)

type RateInfo struct {
	LastRequest time.Time
	Count       int
}

type RateLimiter struct {
	Data  map[string]*RateInfo
	Limit int
	TTL   time.Duration
	Mutex sync.Mutex
}

func NewRateLimiter(cfg Config) *RateLimiter {
	return &RateLimiter{
		Limit: cfg.Limit,
		TTL:   cfg.TTL,
		Data:  make(map[string]*RateInfo),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	info, exists := rl.Data[key]
	now := time.Now()

	if !exists {
		rl.Data[key] = &RateInfo{LastRequest: now, Count: 1}
		go rl.expireKey(key)
		return true
	}

	if now.Sub(info.LastRequest) > rl.TTL {
		info.LastRequest = now
		info.Count = 1
		return true
	}

	if info.Count < rl.Limit {
		info.Count++
		return true
	}

	return false
}

func (rl *RateLimiter) expireKey(key string) {
	time.Sleep(rl.TTL)
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()
	delete(rl.Data, key)
}
