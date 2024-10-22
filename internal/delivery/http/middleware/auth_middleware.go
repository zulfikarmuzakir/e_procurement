package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zulfikarmuzakir/e_procurement/pkg/auth"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

const RoleKey ContextKey = "role"

func JWTAuth(jwtAuth *auth.JWTAuth) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			claims, err := jwtAuth.ValidateToken(bearerToken[1], true)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}
