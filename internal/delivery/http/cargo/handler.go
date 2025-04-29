package cargo

import (
	"log"
	"net/http"
	"strings"
	"test-project/internal/domain/auth"
	cargoDomain "test-project/internal/domain/cargo"
	"test-project/internal/domain/user"
	"test-project/internal/middleware"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc        usecase.CargoUsecase
	deps      *auth.Deps
	validator *validator.Validator
}

func NewHandler(uc usecase.CargoUsecase, deps *auth.Deps, v *validator.Validator) *Handler {
	return &Handler{
		uc:        uc,
		deps:      deps,
		validator: v,
	}
}

func RegisterCargoRoute(r *mux.Router, deps *auth.Deps) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	cargoRepo := cargoDomain.NewPostgresCargoRepo(deps.DB)
	svc := usecase.NewCargoUsecase(cargoRepo, v)
	h := NewHandler(svc, deps, v)

	r.Handle("/cargo", middleware.JwtMiddleware(deps, h.Create)).Methods(http.MethodPost)
	r.Handle("/cargo", middleware.JwtMiddleware(deps, h.GET)).Methods(http.MethodGet)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.PATH)).Methods(http.MethodPatch)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.GETByID)).Methods(http.MethodGet)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.DELETE)).Methods(http.MethodDelete)
}

// Create handles the creation of a new cargo via form-data
// @Summary Create a new cargo
// @Description Creates a new cargo with the provided details via form-data
// @Tags cargo
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param cargoNumber        formData string  true  "Номер груза"
// @Param date               formData string  false "Дата (RFC3339), например 2025-04-30T08:00:00Z"
// @Param loadUnloadDate     formData string  false "Дата погрузки/разгрузки (2025-04-30T08:00:00Z)"
// @Param driver             formData string  true  "Водитель"
// @Param transportationInfo formData string  true  "Информация о перевозке"
// @Param payoutAmount       formData number  false "Сумма выплаты, например 12345.67"
// @Param payoutDate         formData string  false "Дата выплаты (RFC3339)"
// @Param paymentStatus      formData string  false "Статус оплаты"
// @Param payoutTerms        formData string  false "Условия выплаты"
// @Param truckId            formData string  true  "ID машины (c8169351-f6d8-4058-af4a-8ead3363fd92)"
// @Success 201 {object} cargo.CreateResponse "Груз успешно создан"
// @Failure 400 {object} cargo.ErrorResponse  "Ошибки валидации или неверный формат данных"
// @Failure 500 {object} cargo.ErrorResponse  "Внутренняя ошибка сервера"
// @Router /cargo [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.deps.Logger)
		return
	}

	if role != user.RoleSuperAdmin && role != user.RoleEditor {
		utils.JSON(w, http.StatusUnauthorized, "Недостаточно прав. Суперадминистраторы и Редакторы могут создавать груз", nil, h.deps.Logger)
		return
	}
	var c cargoDomain.Cargo
	if err := utils.ParseFormData(r, &c); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Не удалось распарсить форму: "+err.Error(), nil, h.deps.Logger)
		return
	}

	if errs := h.validator.Validate(c); len(errs) > 0 {
		utils.JSON(w, http.StatusBadRequest, "Ошибки валидации: "+strings.Join(errs, "; "), nil, h.deps.Logger)
		return
	}

	created, err := h.uc.CreateCargo(c)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Груз успешно создан", created, h.deps.Logger)
}

// GET retrieves a list of all cargos
// @Summary List all cargos
// @Description Retrieves a list of all cargos in the system
// @Tags cargo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} cargo.ListResponse "List of cargos"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargo [get]
func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	cargos, err := h.uc.ListGargos()

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Список всех грузов", cargos, h.deps.Logger)
}

// GETByID retrieves a cargo by ID
// @Summary Get a cargo by ID
// @Description Retrieves a cargo by its ID
// @Tags cargo
// @Accept json
// @Produce json
// @Param id path string true "Cargo ID"
// @Security BearerAuth
// @Success 200 {object} cargo.GetResponse "Cargo found"
// @Failure 400 {object} cargo.ErrorResponse "Invalid ID"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargo/{id} [get]
func (h *Handler) GETByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	cargo, err := h.uc.GetCargo(id)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Данные о грузе", cargo, h.deps.Logger)
}

// PATH updates a cargo by ID
// @Summary Update a cargo by ID
// @Description Updates a cargo by its ID
// @Tags cargo
// @Accept json
// @Produce json
// @Param id path string true "Cargo ID"
// @Security BearerAuth
// @Param cargo body cargo.CreateRequest true "Cargo object to be updated"
// @Success 200 {object} cargo.GetResponse "Cargo updated"
// @Failure 400 {object} cargo.ErrorResponse "Invalid ID"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargo/{id} [patch]
func (h *Handler) PATH(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.deps.Logger)
		return
	}

	if role != user.RoleSuperAdmin && role != user.RoleEditor {
		utils.JSON(w, http.StatusUnauthorized, "Недостаточно прав. Суперадминистраторы и Редакторы могут обновлять груз", nil, h.deps.Logger)
		return
	}

	id := mux.Vars(r)["id"]

	var updateCargo cargoDomain.UpdateCargoInput

	if err := utils.ParseFormData(r, &updateCargo); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Не удалось распарсить форму: "+err.Error(), nil, h.deps.Logger)
		return
	}

	if errs := h.validator.Validate(updateCargo); len(errs) > 0 {
		utils.JSON(w, http.StatusBadRequest, "Ошибки валидации: "+strings.Join(errs, "; "), nil, h.deps.Logger)
		return
	}

	cargo, err := h.uc.PatchCargo(updateCargo, id)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Груз успешно обновлен", cargo, h.deps.Logger)
}

// DELETE deletes a cargo by ID
// @Summary Delete a cargo by ID
// @Description Deletes a cargo by its ID
// @Tags cargo
// @Accept json
// @Produce json
// @Param id path string true "Cargo ID"
// @Security BearerAuth
// @Success 200 {object} cargo.DeleteResponse "Cargo deleted"
// @Failure 400 {object} cargo.ErrorResponse "Invalid ID"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargo/{id} [delete]
func (h *Handler) DELETE(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.deps.Logger)
		return
	}

	if role != user.RoleSuperAdmin && role != user.RoleEditor {
		utils.JSON(w, http.StatusUnauthorized, "Недостаточно прав. Суперадминистраторы и Редакторы могут удалять груз", nil, h.deps.Logger)
		return
	}

	id := mux.Vars(r)["id"]

	err = h.uc.DeleteCargo(id)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Груз с id= "+id+" успешно удален", nil, h.deps.Logger)
}
