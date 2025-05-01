package user

import "time"

type Role string

const (
	RoleUser       Role = "USER"
	RoleEditor     Role = "EDITOR"
	RoleSuperAdmin Role = "SUPERADMIN"
)

var AllRoles = []Role{
	RoleUser, RoleEditor, RoleSuperAdmin,
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateUser struct {
	Username string `json:"username" validate:"required,min=3"`
	Role     string `json:"role" validate:"required,oneof=USER EDITOR SUPERADMIN"`
}

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id string) (User, error)
	Create(user User) (User, error)
	Delete(id string) error
	FindByEmail(email string) (User, error)
	Update(id string, user UpdateUser) error
}

type ListResponse struct {
	Message string `json:"message" example:"Список пользователей"`
	Data    []User `json:"data"`
}

type GetResponse struct {
	Message string `json:"message" example:"Пользователь"`
	Data    User   `json:"data"`
}

type UpdateRequest struct {
	Username string `json:"username" validate:"required,min=3" example:"username"`
	Role     string `json:"role" validate:"required,oneof=USER EDITOR SUPERADMIN" example:"USER"`
}

type UpdateResponse struct {
	Message string `json:"message" example:"Пользователь"`
	Data    User   `json:"data"`
}

type CreateResponse struct {
	Message string `json:"message" example:"Пользователь успешно создан"`
	Data    User   `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Невалидный формат JSON"`
	Data    interface{} `json:"data"`
}

type DeleteResponse struct {
	Message string      `json:"message" example:"Пользователь успешно удалён"`
	Data    interface{} `json:"data"`
}
