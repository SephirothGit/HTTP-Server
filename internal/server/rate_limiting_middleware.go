package server

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	tokens     int
	lastRefill time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	rate     int
	interval time.Duration
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*client),
		rate:     rate,
		interval: interval,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		rl.mu.Lock()
		defer rl.mu.Unlock()

		c, exists := rl.clients[ip]
		now := time.Now()

		if !exists {
			rl.clients[ip] = &client{
				tokens:     rl.rate - 1,
				lastRefill: now,
			}
			next.ServeHTTP(w, r)
			return
		}

		// Refill tokens
		if now.Sub(c.lastRefill) > rl.interval {
			c.tokens = rl.rate
			c.lastRefill = now
		}

		if c.tokens <= 0 {
			w.WriteHeader(http.StatusTooManyRequests) // Error 429
			w.Write([]byte("rate limit exceeded"))
			return
		}

		c.tokens--
		next.ServeHTTP(w, r)
	})
}
