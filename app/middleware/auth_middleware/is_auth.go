package auth_middleware

import (
	"authentication_backend/utils/jwt"
	"authentication_backend/utils/log"
	"context"
	"net/http"
)

type contextKey string

const userIdKey contextKey = "userId"

func IsAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := jwt.Auth(w, r)

		ctx := context.WithValue(r.Context(), userIdKey, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func GetUserId(ctx context.Context) string {

	userId, ok := ctx.Value(userIdKey).(string)

	if !ok {
		log.Info("No user id from context")
		return ""
	}

	return userId
}
