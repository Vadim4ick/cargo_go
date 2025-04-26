package router

import (
	"test-project/internal/delivery/http/auth"
	"test-project/internal/delivery/http/cargo"
	"test-project/internal/delivery/http/truck"
	"test-project/internal/delivery/http/user"
	"test-project/internal/redis"
	"test-project/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup(pool *pgxpool.Pool, logger *zap.Logger, jwtService *usecase.JwtUsecase, redisService *redis.Client) *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()
	subrouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	user.RegisterUserRoutes(subrouter, pool, logger)
	truck.RegisterUserRoutes(subrouter, pool, logger)
	cargo.RegisterCargoRoute(subrouter, pool, logger)
	auth.RegisterCargoRoute(subrouter, pool, logger, jwtService, redisService)

	return subrouter
}
