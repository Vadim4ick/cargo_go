package router

import (
	user "test-project/internal/delivery/http"
	"test-project/internal/repository"
	"test-project/internal/usecase"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	userRepo := repository.NewInMemoryUserRepo()
	svc := usecase.NewUserUsecase(userRepo)
	h := user.NewHandler(svc)

	r.HandleFunc("/users", h.List).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
}
