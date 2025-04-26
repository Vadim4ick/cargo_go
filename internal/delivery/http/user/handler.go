package user

import (
	"log"
	"net/http"
	"test-project/internal/domain/user"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler struct {
	uc     usecase.UserUsecase
	logger *zap.Logger
}

func NewHandler(uc usecase.UserUsecase, logger *zap.Logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, logger *zap.Logger) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	userRepo := user.NewPostgresUserRepo(db)
	svc := usecase.NewUserUsecase(userRepo, v)
	h := NewHandler(svc, logger)

	r.HandleFunc("/users", h.List).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
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
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список пользователей", users, h.logger)
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
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil, h.logger)
		return
	}

	user, err := h.uc.GetUser(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь", user, h.logger)
}
