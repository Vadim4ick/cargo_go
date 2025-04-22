package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"test-project/internal/domain"
	"test-project/internal/repository"
	"test-project/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	uc usecase.UserUsecase
}

func NewHandler(uc usecase.UserUsecase) *Handler {
	return &Handler{uc: uc}
}

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool) {
	userRepo := repository.NewPostgresUserRepo(db)
	svc := usecase.NewUserUsecase(userRepo)
	h := NewHandler(svc)

	r.HandleFunc("/users", h.List).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
	r.HandleFunc("/users", h.Create).Methods("POST")
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.uc.CreateUser(user)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
