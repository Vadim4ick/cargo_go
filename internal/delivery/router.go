package router

import (
	"test-project/internal/delivery/http/auth"
	"test-project/internal/delivery/http/cargo"
	"test-project/internal/delivery/http/truck"
	"test-project/internal/delivery/http/user"
	authDomain "test-project/internal/domain/auth"
	userDomain "test-project/internal/domain/user"
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

	userRepo := userDomain.NewPostgresUserRepo(pool)
	authSvc := usecase.NewService(userRepo, jwtService, redisService)

	deps := &authDomain.Deps{
		Logger:      logger,
		JwtService:  jwtService,
		AuthService: authSvc,
		Redis:       redisService,
		DB:          pool,
	}

	user.RegisterUserRoutes(subrouter, deps)
	truck.RegisterUserRoutes(subrouter, deps)
	cargo.RegisterCargoRoute(subrouter, deps)
	auth.RegisterCargoRoute(subrouter, deps)

	return subrouter
}
