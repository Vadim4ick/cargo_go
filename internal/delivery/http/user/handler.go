package user

import (
	"encoding/json"
	"log"
	"net/http"
	userDomain "test-project/internal/domain/user"
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

	userRepo := userDomain.NewPostgresUserRepo(db)
	svc := usecase.NewUserUsecase(userRepo, v)
	h := NewHandler(svc)

	r.HandleFunc("/users", h.List).Methods("GET")
	r.HandleFunc("/users/{id}", h.Get).Methods("GET")
	r.HandleFunc("/users", h.Create).Methods("POST")
}

// @Summary      List users
// @Description  Возвращает всех зарегистрированных пользователей
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.ListUsers()
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список пользователей", users)
}

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var user userDomain.User

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
