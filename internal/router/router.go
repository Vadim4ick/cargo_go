package router

import (
	"github.com/gorilla/mux"
)

func Setup() *mux.Router {
	r := mux.NewRouter()

	subrouter := r.PathPrefix("/api/v1").Subrouter()

	RegisterUserRoutes(subrouter)

	return subrouter
}
