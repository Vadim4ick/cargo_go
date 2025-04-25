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
	FindByID(id string) (User, error)
	Create(user User) (User, error)
}
