package cargo

import (
	"encoding/json"
	"log"
	"net/http"
	cargoDomain "test-project/internal/domain/cargo"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

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
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil)
		return
	}

	cargo, err := h.uc.CreateCargo(cargo)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusCreated, "Груз успешно создан", cargo)
}

func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	cargos, err := h.uc.ListGargos()

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusOK, "Список всех грузов", cargos)
}

func (h *Handler) GETByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseNumber(mux.Vars(r)["id"])

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, "Некорректный id", nil)
		return
	}

	cargo, err := h.uc.GetCargo(id)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusOK, "Данные о грузе", cargo)
}
