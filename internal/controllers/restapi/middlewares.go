package restapi

import (
	"context"
	"fmt"
	"github.com/Arh0rn/absoluteCinema/internal/controllers/restapi/controllers"
	"github.com/Arh0rn/absoluteCinema/pkg"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := pkg.GetTokenFromRequest(r)
		if err != nil {
			slog.Error(err.Error(),
				"architecture level", "middleware",
			)
			controllers.HandleError(w, models.ErrInvalidToken, http.StatusUnauthorized)
			return
		}
		id, err := pkg.ParseToken(token, []byte("secret"))
		if err != nil {
			slog.Error(err.Error(),
				"architecture level", "middleware",
			)
			controllers.HandleError(w, models.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "id", id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
