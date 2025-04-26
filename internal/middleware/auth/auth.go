package auth

import (
	"context"
	"net/http"
	"strings"
	"test-project/internal/usecase"
	"test-project/utils"

	"go.uber.org/zap"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func JwtMiddleware(svc usecase.AuthUsecase, jwtSvc *usecase.JwtUsecase, logger *zap.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 {
			utils.JSON(w, http.StatusUnauthorized, "missing token", nil, logger)
			return
		}
		uid, err := jwtSvc.ValidateAccess(parts[1])
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, "invalid token", nil, logger)
			return
		}
		// помечаем онлайн
		svc.TouchOnline(uid)
		// передаём в ctx
		ctx := context.WithValue(r.Context(), UserIDKey, uid)

		next(w, r.WithContext(ctx))
	})
}
