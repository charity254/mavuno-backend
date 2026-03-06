package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
    "github.com/mavuno/mavuno-backend/internal/utils"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextUserRole contextKey = "role"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.Error(w, http.StatusUnauthorized, "authorization header is required")
				return
			}

			//Split "Bearer <token>" into two parts
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.Error(w, http.StatusUnauthorized, "invalid authorization header format")
				return 
			}

			tokenString := parts[1]

			//Parse and validate the JWT token using the secret
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				utils.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return 
			}
			//EXtract the claims (user_id and role) from the token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				utils.Error(w, http.StatusUnauthorized, "invalid token claims")
				return 
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims["user_id"])
			ctx = context.WithValue(ctx, ContextUserRole, claims["role"])


			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func RequiredRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextUserRole).(string)
			if !ok || userRole != role {
				utils.Error(w, http.StatusForbidden, "you do not have permission to access this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}