package restapi

import (
	"log/slog"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request",
			"method", r.Method,
			"url", r.URL,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent())
		//log.Println(r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}

func AddMiddlewares(router http.Handler) http.Handler {
	return loggingMiddleware(router) // add wrapping here
}
