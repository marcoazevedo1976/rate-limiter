package limiter

import (
	"net"
	"net/http"
	"strings"
)

func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("API_KEY")
		isToken := token != ""
		var key string

		if isToken {
			key = "token:" + token
		} else {
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}
			// Em ambientes com proxy
			if fwd := req.Header.Get("X-Forwarded-For"); fwd != "" {
				parts := strings.Split(fwd, ",")
				ip = strings.TrimSpace(parts[0])
			}
			key = "ip:" + ip
		}

		allowed, status, message := r.Allow(key, isToken)
		if !allowed {
			http.Error(w, message, status)
			return
		}

		next.ServeHTTP(w, req)
	})
}
