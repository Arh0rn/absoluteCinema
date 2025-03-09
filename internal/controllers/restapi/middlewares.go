package restapi

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func createMiddlewareStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		end := time.Since(start)

		slog.Info("Incoming request",
			"method", r.Method,
			"url", r.URL,
			"remote_addr", r.RemoteAddr,
			"duration", fmt.Sprintf("%v miliseconds", end.Milliseconds()),
			"user_agent", r.UserAgent(),
		)
	})
}
