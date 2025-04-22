package router

import (
	user "test-project/internal/delivery/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(pool *pgxpool.Pool) *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()

	user.RegisterUserRoutes(subrouter, pool)

	return subrouter
}
