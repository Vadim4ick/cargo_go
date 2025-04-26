package router

import (
	"test-project/internal/delivery/http/auth"
	"test-project/internal/delivery/http/cargo"
	"test-project/internal/delivery/http/truck"
	"test-project/internal/delivery/http/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup(pool *pgxpool.Pool, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()
	subrouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	user.RegisterUserRoutes(subrouter, pool, logger)
	truck.RegisterUserRoutes(subrouter, pool, logger)
	cargo.RegisterCargoRoute(subrouter, pool, logger)
	auth.RegisterCargoRoute(subrouter, pool, logger)

	return subrouter
}
