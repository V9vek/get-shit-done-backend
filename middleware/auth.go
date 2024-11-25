package middleware

import (
	"fmt"
	"get-shit-done/service"
	"net/http"
	"strings"
)

func ValidateAccessToken(jwtAuth *service.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if auth == "" {
				http.Error(w, "missing or malformed token", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(auth, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "missing or malformed token", http.StatusUnauthorized)
				return
			}

			token := headerParts[1]
			isValid, err := jwtAuth.IsAccessTokenValid(token)
			if !isValid || err != nil {
				http.Error(w, fmt.Sprintf("invalid or expired token %v", err), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
