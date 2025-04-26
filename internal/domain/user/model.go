package user

import "time"

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id int) (User, error)
	Create(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type ListResponse struct {
	Message string `json:"message" example:"Список пользователей"`
	Data    []User `json:"data"`
}

type GetResponse struct {
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
