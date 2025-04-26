package user

import (
	"encoding/json"
	"log"
	"net/http"
	"test-project/internal/domain/user"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

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
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	userRepo := user.NewPostgresUserRepo(db)
	svc := usecase.NewUserUsecase(userRepo, v)
	h := NewHandler(svc)

	r.HandleFunc("/users", h.List).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
	r.HandleFunc("/users", h.Create).Methods("POST")
}

// List retrieves a list of all users
// @Summary List all users
// @Description Retrieves a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} user.ListResponse "List of users"
// @Failure 404 {object} user.ErrorResponse "Users not found"
// @Router /users [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список пользователей", users)
}

// Get retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieves a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 201 {object} user.GetResponse "User found"
// @Failure 400 {object} user.ErrorResponse "Invalid ID"
// @Failure 404 {object} user.ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseNumber(mux.Vars(r)["id"])
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil)
		return
	}

	user, err := h.uc.GetUser(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь", user)
}

// Create handles the creation of a new user
// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.User true "User object to be created"
// @Success 201 {object} user.CreateResponse "User successfully created"
// @Failure 400 {object} user.ErrorResponse "Invalid JSON format or creation error"
// @Router /users [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user user.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil)
		return
	}

	user, err := h.uc.CreateUser(user)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь успешно создан", user)
}
