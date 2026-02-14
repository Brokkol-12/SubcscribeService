package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request completed",
			"Method", r.Method,
			"Path", r.URL.Path,
			"Duration", time.Since(start).String(),
			"Remote_addr", r.RemoteAddr)
	})
}
