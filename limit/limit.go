package limit

import (
	"net/http"
	"runtime"
)

var limiter = NewLimiter(uint64(runtime.NumCPU() * 200000))

func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
