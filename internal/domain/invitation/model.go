package invitation

import "time"

type Invitation struct {
	ID int `json:"id"`

	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
	Used  bool   `json:"used" default:"false"`

	CreatedAt time.Time `json:"createdAt"`
}

type InvitationRepository interface {
	Create(invitation Invitation) (Invitation, error)
}

type CreateResponse struct {
	Message string     `json:"message" example:"Приглашение успешно создано"`
	Data    Invitation `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Невалидный формат JSON"`
	Data    interface{} `json:"data"`
}
