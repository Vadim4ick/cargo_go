package cargo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	cargoDomain "test-project/internal/domain/cargo"
	"test-project/internal/usecase"
	"test-project/internal/validator"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	uc usecase.CargoUsecase
}

func NewHandler(uc usecase.CargoUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func RegisterCargoRoute(r *mux.Router, db *pgxpool.Pool) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	cargoRepo := cargoDomain.NewPostgresCargoRepo(db)
	svc := usecase.NewCargoUsecase(cargoRepo, v)
	h := NewHandler(svc)

	r.HandleFunc("/cargos", h.Create).Methods("POST")
	r.HandleFunc("/cargos", h.GET).Methods("GET")
	r.HandleFunc("/cargos/{id}", h.GETByID).Methods("GET")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var cargo cargoDomain.Cargo

	if err := json.NewDecoder(r.Body).Decode(&cargo); err != nil {
		http.Error(w, "Невалидный формат JSON", http.StatusBadRequest)
		return
	}

	cargo, err := h.uc.CreateCargo(cargo)

	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cargo)
}

func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	cargos, err := h.uc.ListGargos()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cargos)
}

func (h *Handler) GETByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	// Теперь id — это int
	cargo, err := h.uc.GetCargo(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cargo)
}
