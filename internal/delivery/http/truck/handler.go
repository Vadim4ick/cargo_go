package truck

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test-project/internal/domain/auth"
	truckDomain "test-project/internal/domain/truck"
	"test-project/internal/domain/user"
	"test-project/internal/middleware"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc   usecase.TruckUsecase
	deps *auth.Deps
}

func NewHandler(uc usecase.TruckUsecase, deps *auth.Deps) *Handler {
	return &Handler{uc: uc, deps: deps}
}
func RegisterUserRoutes(r *mux.Router, deps *auth.Deps) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	truckRepo := truckDomain.NewPostgresTruckRepo(deps.DB)
	svc := usecase.NewTruckUsecase(truckRepo, v)
	h := NewHandler(svc, deps)

	r.Handle("/truck", middleware.JwtMiddleware(deps, h.Create)).Methods(http.MethodPost)
	r.Handle("/truck", middleware.JwtMiddleware(deps, h.GET)).Methods(http.MethodGet)
	r.Handle("/truck/{id}", middleware.JwtMiddleware(deps, h.GETById)).Methods(http.MethodGet)
	r.Handle("/truck/{id}/cargos", middleware.JwtMiddleware(deps, h.GETCargos)).Methods(http.MethodGet)
}

// Create handles the creation of a new truck
// @Summary Create a new truck
// @Description Creates a new truck with the provided details
// @Tags truck
// @Accept json
// @Produce json
// @Param truck body truck.CreateRequest true "Truck object to be created"
// @Security BearerAuth
// @Success 201 {object} cargo.CreateResponse "Cargo successfully created"
// @Failure 400 {object} cargo.ErrorResponse "Invalid JSON format"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /truck [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	role, err := middleware.GetUserRole(r.Context())

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, err.Error(), nil, h.deps.Logger)
		return
	}

	if role != user.RoleSuperAdmin && role != user.RoleEditor {
		utils.JSON(w, http.StatusUnauthorized, "Недостаточно прав. Суперадминистраторы и редакторы могут создавать машины", nil, h.deps.Logger)
		return
	}

	var truck truckDomain.Truck

	if err := json.NewDecoder(r.Body).Decode(&truck); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil, h.deps.Logger)
		return
	}

	truck, err = h.uc.CreateTruck(truck)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Машина успешно создана", truck, h.deps.Logger)
}

// GET retrieves a list of all trucks
// @Summary List all trucks
// @Description Retrieves a list of all trucks in the system
// @Tags truck
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} truck.ListResponse "List of trucks"
// @Failure 404 {object} truck.ErrorResponse "Trucks not found"
// @Router /truck [get]
func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	trucks, err := h.uc.ListTrucks()
	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список машин", trucks, h.deps.Logger)
}

// GETById retrieves a truck by ID
// @Summary Get a truck by ID
// @Description Retrieves a truck by their ID
// @Tags truck
// @Accept json
// @Produce json
// @Param id path string true "Truck ID"
// @Security BearerAuth
// @Success 201 {object} truck.GetResponse "Truck found"
// @Failure 400 {object} truck.ErrorResponse "Invalid ID"
// @Failure 404 {object} truck.ErrorResponse "Truck not found"
// @Router /truck/{id} [get]
func (h *Handler) GETById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	truck, err := h.uc.GetTruck(id)

	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Машина", truck, h.deps.Logger)
}

// GETCargos retrieves a list of cargos by truck ID
// @Summary Get a list of cargos by truck ID
// @Description Retrieves a list of cargos by their truck ID
// @Tags truck
// @Accept json
// @Produce json
// @Param id path string true "Truck ID"
// @Param limit query int false "Limit number of cargos per page"
// @Param page query int false "Page number for pagination"
// @Security BearerAuth
// @Success 201 {object} cargo.ListResponse "List of cargos"
// @Failure 400 {object} cargo.ErrorResponse "Invalid ID"
// @Failure 404 {object} cargo.ErrorResponse "Truck not found"
// @Router /truck/{id}/cargos [get]
func (h *Handler) GETCargos(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Пагинация
	query := r.URL.Query()
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	// Можешь распарсить строки в int
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	cargos, err := h.uc.GetTruckCargos(id, limit, page)

	if err != nil {
		utils.JSON(w, http.StatusNotFound, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Список грузов", cargos, h.deps.Logger)
}
