package user

import (
	"log"
	"net/http"
	"strconv"
	"test-project/internal/domain/auth"
	"test-project/internal/domain/user"
	"test-project/internal/middleware"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc   usecase.UserUsecase
	deps *auth.Deps
}

func NewHandler(uc usecase.UserUsecase, deps *auth.Deps) *Handler {
	return &Handler{uc: uc, deps: deps}
}

func RegisterUserRoutes(r *mux.Router, deps *auth.Deps) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	userRepo := user.NewPostgresUserRepo(deps.DB)
	svc := usecase.NewUserUsecase(userRepo, v)
	h := NewHandler(svc, deps)

	r.Handle("/users", middleware.JwtMiddleware(deps, h.List)).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.Delete).Methods(http.MethodDelete)
}

// List retrieves a list of all users
// @Summary List all users
// @Description Retrieves a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} user.ListResponse "List of users"
// @Failure 404 {object} user.ErrorResponse "Users not found"
// @Router /users [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список пользователей", users, h.deps.Logger)
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
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil, h.deps.Logger)
		return
	}

	user, err := h.uc.GetUser(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь", user, h.deps.Logger)
}

// Delete deletes a user by ID
// @Summary Delete a user by ID
// @Description Deletes a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 201 {object} user.DeleteResponse "User deleted"
// @Failure 400 {object} user.ErrorResponse "Invalid ID"
// @Failure 404 {object} user.ErrorResponse "User not found"
// @Router /users/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseNumber(mux.Vars(r)["id"])
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil, h.deps.Logger)
		return
	}

	err = h.uc.DeleteUser(id)
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Пользователь с id= "+strconv.Itoa(id)+" успешно удален", nil, h.deps.Logger)
}
