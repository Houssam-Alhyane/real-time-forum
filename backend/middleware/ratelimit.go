package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	count int
	reset time.Time
}

// RateLimiter stores per-client request counters for a fixed time window.
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]bucket
	limit   int
	window  time.Duration
}

// NewRateLimiter creates a limiter with safe defaults when inputs are invalid.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &RateLimiter{
		clients: make(map[string]bucket),
		limit:   limit,
		window:  window,
	}
}

// Middleware blocks requests when a client exceeds the allowed requests per window.
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip := clientIP(r.RemoteAddr)

		rl.mu.Lock()
		// Start a new window once the previous one expires.
		b := rl.clients[ip]
		if now.After(b.reset) {
			b = bucket{reset: now.Add(rl.window)}
		}

		if b.count >= rl.limit {
			rl.mu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error":"too many requests"}`))
			return
		}

		b.count++
		rl.clients[ip] = b
		rl.mu.Unlock()

		next(w, r)
	}
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}
	return remoteAddr
}
