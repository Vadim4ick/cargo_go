package usecase

import (
	"errors"
	"strings"
	invitationDomain "test-project/internal/domain/invitation"
	"test-project/internal/validator"

	"github.com/jackc/pgx/v5/pgconn"
)

type InvitationUsecase interface {
	CreateInvitation(invitation invitationDomain.Invitation) (invitationDomain.Invitation, error)
}

type invitationUsecase struct {
	repo      invitationDomain.InvitationRepository
	validator *validator.Validator
}

func NewInvitationService(r invitationDomain.InvitationRepository, validator *validator.Validator) InvitationUsecase {
	return &invitationUsecase{repo: r, validator: validator}
}

func (u *invitationUsecase) CreateInvitation(invitation invitationDomain.Invitation) (invitationDomain.Invitation, error) {
	if errs := u.validator.Validate(invitation); len(errs) > 0 {
		return invitationDomain.Invitation{}, errors.New(strings.Join(errs, "; "))
	}

	invitation, err := u.repo.Create(invitation)
	if err != nil {
		// обработка ошибки дубликата
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return invitationDomain.Invitation{}, errors.New("Приглашение для этого email уже отправлено")
		}
		return invitationDomain.Invitation{}, err
	}

	return invitation, nil
}
