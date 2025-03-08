package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}
func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exist := rl.clients[ip]
	rl.RUnlock()
	if !exist || count < rl.limit {
		rl.Lock()
		if !exist {
			//this allow reset if the client with ip reach their limit on time frame
			go rl.resetCount(ip)
		}
		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}
	return false, rl.window
}
func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
