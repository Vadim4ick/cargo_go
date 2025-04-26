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

// Create handles the creation of a new cargo
// @Summary Create a new cargo
// @Description Creates a new cargo with the provided details
// @Tags cargos
// @Accept json
// @Produce json
// @Param cargo body cargo.Cargo true "Cargo object to be created"
// @Success 201 {object} cargo.CreateResponse "Cargo successfully created"
// @Failure 400 {object} cargo.ErrorResponse "Invalid JSON format"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargos [post]
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

// GET retrieves a list of all cargos
// @Summary List all cargos
// @Description Retrieves a list of all cargos in the system
// @Tags cargos
// @Accept json
// @Produce json
// @Success 200 {object} cargo.ListResponse "List of cargos"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargos [get]
func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	cargos, err := h.uc.ListGargos()

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(w, http.StatusOK, "Список всех грузов", cargos)
}

// GETByID retrieves a cargo by ID
// @Summary Get a cargo by ID
// @Description Retrieves a cargo by its ID
// @Tags cargos
// @Accept json
// @Produce json
// @Param id path int true "Cargo ID"
// @Success 200 {object} cargo.GetResponse "Cargo found"
// @Failure 400 {object} cargo.ErrorResponse "Invalid ID"
// @Failure 500 {object} cargo.ErrorResponse "Internal server error"
// @Router /cargos/{id} [get]
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
