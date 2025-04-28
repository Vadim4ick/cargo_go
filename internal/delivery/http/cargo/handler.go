package cargo

import (
	"encoding/json"
	"log"
	"net/http"
	"test-project/internal/domain/auth"
	cargoDomain "test-project/internal/domain/cargo"
	"test-project/internal/middleware"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc   usecase.CargoUsecase
	deps *auth.Deps
}

func NewHandler(uc usecase.CargoUsecase, deps *auth.Deps) *Handler {
	return &Handler{
		uc:   uc,
		deps: deps,
	}
}

func RegisterCargoRoute(r *mux.Router, deps *auth.Deps) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	cargoRepo := cargoDomain.NewPostgresCargoRepo(deps.DB)
	svc := usecase.NewCargoUsecase(cargoRepo, v)
	h := NewHandler(svc, deps)

	r.Handle("/cargo", middleware.JwtMiddleware(deps, h.Create)).Methods(http.MethodPost)
	r.Handle("/cargo", middleware.JwtMiddleware(deps, h.GET)).Methods(http.MethodGet)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.PATH)).Methods(http.MethodPatch)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.GETByID)).Methods(http.MethodGet)
	r.Handle("/cargo/{id}", middleware.JwtMiddleware(deps, h.DELETE)).Methods(http.MethodDelete)
}

// Create handles the creation of a new cargo
// @Summary Create a new cargo
// @Description Creates a new cargo with the provided details
// @Tags cargo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cargo body cargo.CreateRequest true "Cargo object to be created"
// @Success 201 {object} cargo.CreateResponse "Cargo successfully created"
// @Failure 400 {object} cargo.ErrorResponse "Invalid JSON format"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargo [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var cargo cargoDomain.Cargo

	if err := json.NewDecoder(r.Body).Decode(&cargo); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil, h.deps.Logger)
		return
	}

	cargo, err := h.uc.CreateCargo(cargo)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Груз успешно создан", cargo, h.deps.Logger)
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
	id := mux.Vars(r)["id"]

	var updateCargo cargoDomain.UpdateCargoInput

	if err := json.NewDecoder(r.Body).Decode(&updateCargo); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil, h.deps.Logger)
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
	id := mux.Vars(r)["id"]

	err := h.uc.DeleteCargo(id)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusOK, "Груз с id= "+id+" успешно удален", nil, h.deps.Logger)
}
