package truck

import (
	"encoding/json"
	"log"
	"net/http"
	truckDomain "test-project/internal/domain/truck"
	"test-project/internal/usecase"
	"test-project/internal/validator"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	uc usecase.TruckUsecase
}

func NewHandler(uc usecase.TruckUsecase) *Handler {
	return &Handler{uc: uc}
}
func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	truckRepo := truckDomain.NewPostgresTruckRepo(db)
	svc := usecase.NewTruckUsecase(truckRepo, v)
	h := NewHandler(svc)

	r.HandleFunc("/trucks", h.Create).Methods("POST")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var truck truckDomain.Truck

	if err := json.NewDecoder(r.Body).Decode(&truck); err != nil {
		http.Error(w, "Невалидный формат JSON", http.StatusBadRequest)
		return
	}

	truck, err := h.uc.CreateTruck(truck)

	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(truck)
}
