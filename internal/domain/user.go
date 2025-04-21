package domain

import "time"

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      Role   `json:"role"`
	CreatedAt time.Time
}

type UserRepository interface {
	FindAll() ([]User, error)
	FindByID(id string) (User, error)
	Create(user User) (User, error)
}
