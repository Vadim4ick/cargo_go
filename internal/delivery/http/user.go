package user

import (
	"encoding/json"
	"net/http"
	"test-project/internal/usecase"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc usecase.UserUsecase
}

func NewHandler(uc usecase.UserUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := h.uc.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
