package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

// JWTMiddleware проверяет и валидирует JWT, добавляет userID в context
func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "no token", http.StatusUnauthorized)
				return
			}

			// Убираем "Bearer " и получаем сам токен
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			// Разбираем токен
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Извлекаем userID из токена
			claims := token.Claims.(jwt.MapClaims)
			userID := int(claims["user_id"].(float64))

			// Добавляем userID в контекст
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext извлекает userID из контекста запроса
func UserIDFromContext(ctx context.Context) int {
	return ctx.Value(userIDKey).(int)
}
