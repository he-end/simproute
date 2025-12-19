package ratelimiter

import (
	"net"
	"net/http"
	"time"

	logger "github.com/he-end/simproute/route_logger"
	"go.uber.org/zap"
)

func AddRateLimit(path string, capacity, refill int, interval time.Duration) func(http.Handler) http.Handler {
	rl := NewRateLimiter(capacity, refill, interval)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == path {
				ip, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					ip = r.RemoteAddr
				}
				if !rl.allow(ip) {
					logger.GetLogger().Warn("Rate limit exceeded", zap.String("ip", ip), zap.String("path", r.URL.Path))
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write([]byte("too many requests"))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
