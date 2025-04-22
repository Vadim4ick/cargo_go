package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"test-project/internal/domain"
	"test-project/internal/repository"
	"test-project/internal/usecase"
	"test-project/internal/validator"

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
	userValidator, err := validator.NewUserValidator()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	userRepo := repository.NewPostgresUserRepo(db)
	svc := usecase.NewUserUsecase(userRepo, userValidator)
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
		http.Error(w, "Невалидный формат JSON", http.StatusBadRequest)
		return
	}

	user, err := h.uc.CreateUser(user)
	if err != nil {
		// Ошибки валидации (400), всё остальное (500)
		status := http.StatusInternalServerError
		if err.Error() == "Имя пользователя обязательно" ||
			strings.HasPrefix(err.Error(), "Невалидный") ||
			strings.HasPrefix(err.Error(), "Пароль") {
			status = http.StatusBadRequest
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
