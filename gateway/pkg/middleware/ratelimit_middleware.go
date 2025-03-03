package middleware

import (
	"fmt"
	"github.com/vietquan-37/gateway/pkg/ratelimiter"
	"net/http"
)

const (
	msg = "Too many requests please try later in %s"
)

type RateLimiterMiddleware struct {
	limiter ratelimiter.Limiter
}

func NewRateLimiterMiddleware(limiter ratelimiter.Limiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: limiter,
	}
}
func (m *RateLimiterMiddleware) RateLimitMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allow, retryAfter := m.limiter.Allow(r.RemoteAddr); !allow {
			http.Error(w, fmt.Sprintf(msg, retryAfter.String()), http.StatusTooManyRequests)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
