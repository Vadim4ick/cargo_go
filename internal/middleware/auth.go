package middleware

import (
	"context"
	"net/http"
	"strings"
	"test-project/internal/domain/auth"
	"test-project/utils"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func JwtMiddleware(deps *auth.Deps, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 {
			utils.JSON(w, http.StatusUnauthorized, "missing token", nil, deps.Logger)
			return
		}
		uid, err := deps.JwtService.ValidateAccess(parts[1])
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, "invalid token", nil, deps.Logger)
			return
		}
		// помечаем онлайн
		deps.AuthService.TouchOnline(uid)
		// передаём в ctx
		ctx := context.WithValue(r.Context(), UserIDKey, uid)

		next(w, r.WithContext(ctx))
	})
}
