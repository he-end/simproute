package ratelimiter

import (
	"sync"
	"time"
)

// Simple token bucket per IP
type bucket struct {
	tokens         int
	lastRefillTime time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	capacity int
	refill   int
	interval time.Duration
}

func NewRateLimiter(capacity int, refill int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*bucket),
		capacity: capacity,
		refill:   refill,
		interval: interval,
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[ip]
	if !ok {
		b = &bucket{tokens: rl.capacity, lastRefillTime: time.Now()}
		rl.buckets[ip] = b
	}
	// refill
	now := time.Now()
	elapsed := now.Sub(b.lastRefillTime)
	if elapsed >= rl.interval {
		intervals := int(elapsed / rl.interval)
		b.tokens += intervals * rl.refill
		if b.tokens > rl.capacity {
			b.tokens = rl.capacity
		}
		b.lastRefillTime = now
	}
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}
