package auth

import (
	"test-project/internal/domain/user"
	"test-project/internal/redis"
	"test-project/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	Email       string `json:"email" example:"firulvv@mail.ru"`
	Password    string `json:"password" example:"123456"`
	InviteToken string `json:"inviteToken" example:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"firulvv@mail.ru"`
	Password string `json:"password" example:"123456"`
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

type Deps struct {
	Logger      *zap.Logger
	JwtService  *usecase.JwtUsecase
	AuthService usecase.AuthUsecase
	Redis       *redis.Client
	DB          *pgxpool.Pool
	FileService *usecase.FileService
}
