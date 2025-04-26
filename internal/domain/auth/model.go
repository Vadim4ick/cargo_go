package auth

import (
	"test-project/internal/domain/user"
)

// type Service interface {
// 	Register(email, password string) (user.User, error)
// 	Login(email, password string) (token string, err error)
// 	// помечаем пользователя онлайн (вызыв. на каждом запросе)
// 	TouchOnline(userID string) error
// 	// возвращает список онлайн-ID
// 	OnlineUsers(since time.Duration) ([]string, error)
// 	GetUser(id int) (user.User, error)
// }

type RegisterRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"securepass"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"securepass"`
}

type RegisterResponse struct {
	Message string    `json:"message" example:"Пользователь успешно зарегистрирован"`
	Data    user.User `json:"data"`
}

type LoginResponse struct {
	Message string            `json:"message" example:"Пользователь успешно авторизован"`
	Data    map[string]string `json:"data"`
}

type RefreshResponse struct {
	Message string            `json:"message" example:"Токен успешно обновлён"`
	Data    map[string]string `json:"data"`
}

type OnlineListResponse struct {
	Message string   `json:"message" example:"Список ID-шников онлайн пользователей"`
	Data    []string `json:"data"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Успешный выход из системы"`
}

type ProfileResponse struct {
	Message string    `json:"message" example:"Данные о пользователе"`
	Data    user.User `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Некорректные данные"`
	Data    interface{} `json:"data"`
}
