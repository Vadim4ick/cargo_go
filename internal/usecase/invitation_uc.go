package usecase

import (
	"errors"
	"strings"
	invitationDomain "test-project/internal/domain/invitation"
	"test-project/internal/validator"
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

	return u.repo.Create(invitation)
}
