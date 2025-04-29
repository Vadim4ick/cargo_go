package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"test-project/internal/domain/auth"
	"test-project/internal/domain/user"
	"test-project/utils"
)

type ctxKey string

const (
	UserIDKey   ctxKey = "userID"
	UserRoleKey ctxKey = "userRole"
)

func JwtMiddleware(deps *auth.Deps, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 {
			utils.JSON(w, http.StatusUnauthorized, "missing token", nil, deps.Logger)
			return
		}
		uid, role, err := deps.JwtService.ValidateAccess(parts[1])
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, deps.Logger)
			return
		}

		// помечаем онлайн
		deps.AuthService.TouchOnline(uid)
		// передаём в ctx
		ctx := context.WithValue(r.Context(), UserIDKey, uid)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next(w, r.WithContext(ctx))
	})
}

func GetUserRole(ctx context.Context) (user.Role, error) {
	val := ctx.Value(UserRoleKey)
	if val == nil {
		return "", errors.New("role not found in context")
	}
	role, ok := val.(user.Role)
	if !ok {
		return "", errors.New("invalid role type in context")
	}
	return role, nil
}
