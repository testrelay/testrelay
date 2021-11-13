package httputil

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// LogMiddleware logs the inbound request.
func LogMiddleware(logger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/healthz" {
				logger.Info(r.URL.Path)
			}

			h.ServeHTTP(w, r)
		})
	}
}

// RequireAccessTokenMiddleware requires a valid authorization header is set and matches
// the token provided.
func RequireAccessTokenMiddleware(token string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			passed := r.Header.Get("Authorization")
			if passed != token {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"errors":["invalid access token"]}`))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
