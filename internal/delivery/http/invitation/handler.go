package invitation

import (
	"encoding/json"
	"log"
	"net/http"
	"test-project/internal/domain/auth"
	invitationDomain "test-project/internal/domain/invitation"
	"test-project/internal/usecase"
	"test-project/internal/validator"
	"test-project/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	uc   usecase.InvitationUsecase
	deps *auth.Deps
}

func NewHandler(uc usecase.InvitationUsecase, deps *auth.Deps) *Handler {
	return &Handler{
		uc:   uc,
		deps: deps,
	}
}

func RegisterInvitationRoutes(r *mux.Router, deps *auth.Deps) {
	v, err := validator.New()
	if err != nil {
		log.Fatal("Ошибка инициализации валидатора:", err)
	}

	invitationRepo := invitationDomain.NewPostgresCargoRepo(deps.DB)
	svc := usecase.NewInvitationService(invitationRepo, v)

	h := NewHandler(svc, deps)

	r.HandleFunc("/invitation", h.CREATE).Methods(http.MethodPost)
}

// CREATE handles the creation of a new invitation
// @Summary Create a new invitation
// @Description Creates a new invitation with the provided details
// @Tags invitations
// @Accept json
// @Produce json
// @Param invitation body invitation.Invitation true "Invitation object to be created"
// @Success 201 {object} invitation.CreateResponse "Invitation successfully created"
// @Failure 400 {object} invitation.ErrorResponse "Invalid JSON format"
// @Failure 500 {object} invitation.ErrorResponse "Internal server error"
// @Router /invitation [post]
func (h *Handler) CREATE(w http.ResponseWriter, r *http.Request) {

	var cargo invitationDomain.Invitation

	if err := json.NewDecoder(r.Body).Decode(&cargo); err != nil {
		utils.JSON(w, http.StatusBadRequest, "Невалидный формат JSON", nil, h.deps.Logger)
		return
	}

	cargo, err := h.uc.CreateInvitation(cargo)

	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, err.Error(), nil, h.deps.Logger)
		return
	}

	utils.JSON(w, http.StatusCreated, "Приглашение успешно создано", cargo, h.deps.Logger)
}
