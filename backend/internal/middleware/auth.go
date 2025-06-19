package middleware

import (
	"net/http"
)

// AuthMiddleware is a stub for JWT authentication middleware.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement JWT auth check
		next.ServeHTTP(w, r)
	})
}
