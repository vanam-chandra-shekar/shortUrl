package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func GlobalRateLimiter(rps rate.Limit, burst int) MiddleWare {

	return func(next http.Handler) http.Handler {

		limiter := rate.NewLimiter(rps, burst)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "To many Reauests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
