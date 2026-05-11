package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/mavuno/mavuno-backend/internal/utils"
)

// ipLimiter tracks request counts per IP address
type ipLimiter struct {
	count    int
	lastSeen time.Time
}

var (
	limiters = make(map[string]*ipLimiter)
	mu       sync.Mutex
	maxReqs  = 100              // maximum requests per window
	window   = 1 * time.Minute // time window
)

// RateLimitMiddleware limits each IP to 100 requests per minute.
func RateLimitMiddleware(next http.Handler) http.Handler {
	// Start a background goroutine to clean up old entries every minute
	go cleanupLimiters()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists || time.Since(limiter.lastSeen) > window {
			// New IP or window has expired — reset the counter
			limiters[ip] = &ipLimiter{count: 1, lastSeen: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		limiter.count++
		if limiter.count > maxReqs {
			mu.Unlock()
			utils.Error(w, http.StatusTooManyRequests, "too many requests — please try again later")
			return
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// cleanupLimiters removes IP entries that have not been seen in over a minute.
// Runs every minute in the background to prevent memory leaks.
func cleanupLimiters() {
	for {
		time.Sleep(window)
		mu.Lock()
		for ip, limiter := range limiters {
			if time.Since(limiter.lastSeen) > window {
				delete(limiters, ip)
			}
		}
		mu.Unlock()
	}
}