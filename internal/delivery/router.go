package router

import (
	user "test-project/internal/delivery/http"

	"github.com/gorilla/mux"
)

func Setup() *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()

	user.RegisterUserRoutes(subrouter)

	return subrouter
}
