package middleware

import (
	"fmt"
	"get-shit-done/service"
	"net/http"
)

func ValidateAccessToken(jwtAuth *service.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/*
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
			*/
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, fmt.Sprintf("access token not found: %v", err), http.StatusUnauthorized)
				return
			}

			accessToken := accessTokenCookie.Value

			isValid, err := jwtAuth.IsAccessTokenValid(accessToken)
			if !isValid || err != nil {
				http.Error(w, fmt.Sprintf("invalid or expired token %v", err), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
