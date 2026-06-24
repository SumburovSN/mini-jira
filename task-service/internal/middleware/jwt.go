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
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid claims", http.StatusUnauthorized)
				return
			}

			uidRaw, ok := claims["user_id"]
			if !ok {
				http.Error(w, "user_id missing", http.StatusUnauthorized)
				return
			}

			userID := int(uidRaw.(float64))

			// Добавляем userID в контекст
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext извлекает userID из контекста запроса
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
