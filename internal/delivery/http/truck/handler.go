package truck

import (
	"encoding/json"
	"log"
	"net/http"
	truckDomain "test-project/internal/domain/truck"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Handler struct {
	uc     usecase.TruckUsecase
	logger *zap.Logger
}

func NewHandler(uc usecase.TruckUsecase, logger *zap.Logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}
func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, logger *zap.Logger) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	truckRepo := truckDomain.NewPostgresTruckRepo(db)
	svc := usecase.NewTruckUsecase(truckRepo, v)
	h := NewHandler(svc, logger)

	r.HandleFunc("/trucks", h.Create).Methods("POST")
}

// Create handles the creation of a new truck
// @Summary Create a new truck
// @Description Creates a new truck with the provided details
// @Tags trucks
// @Accept json
// @Produce json
// @Param truck body truck.Truck true "Truck object to be created"
// @Success 201 {object} cargo.CreateResponse "Cargo successfully created"
// @Failure 400 {object} cargo.ErrorResponse "Invalid JSON format"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /trucks [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var truck truckDomain.Truck

	if err := json.NewDecoder(r.Body).Decode(&truck); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil, h.logger)
		return
	}

	truck, err := h.uc.CreateTruck(truck)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Машина успешно создана", truck, h.logger)
}
