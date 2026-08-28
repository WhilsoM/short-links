package custommiddlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"short-links/internal/utils"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	return userID, nil
}

func AuthMiddleware(jwtmanager utils.JWTManagerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				utils.WriteJSONErrorResponse(w, http.StatusUnauthorized, "headers missed")
				return
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHeader, prefix) {
				utils.WriteJSONErrorResponse(w, http.StatusUnauthorized, "doesnt have bearer prefix")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, prefix)

			if tokenString == "" {
				utils.WriteJSONErrorResponse(w, http.StatusUnauthorized, "authorization token is required")
				return
			}

			userID, err := jwtmanager.ParseToken(tokenString)
			if err != nil {
				slog.Info("failed to parse token", "error", err)
				utils.WriteJSONErrorResponse(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
