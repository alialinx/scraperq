package middleware

import (
	"net/http"
	"strings"

	"github.com/alialin/scraperq/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")

		if header == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := auth.ValidateToken(tokenString, jwtSecret)

		claims, ok := token.Claims.(jwt.MapClaims)

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)

	}
}
