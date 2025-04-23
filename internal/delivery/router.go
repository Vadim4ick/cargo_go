package router

import (
	"test-project/internal/delivery/http/truck"
	"test-project/internal/delivery/http/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup(pool *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()
	subrouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	user.RegisterUserRoutes(subrouter, pool)
	truck.RegisterUserRoutes(subrouter, pool)

	return subrouter
}
