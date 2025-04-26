package auth

import (
	"context"
	"net/http"
	"strings"
	"test-project/internal/usecase"
	"test-project/utils"
	"time"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// jwtMiddleware извлекает токен, валидирует, кладёт userID в контекст
// и вызывает TouchOnline
func JwtMiddleware(svc usecase.AuthUsecase, jwtSvc *usecase.JwtUsecase, idle time.Duration, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 {
			utils.JSON(w, http.StatusUnauthorized, "missing token", nil)
			return
		}
		uid, err := jwtSvc.ValidateAccess(parts[1])
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, "invalid token", nil)
			return
		}
		// помечаем онлайн
		svc.TouchOnline(uid)
		// передаём в ctx
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next(w, r.WithContext(ctx))
	})
}
