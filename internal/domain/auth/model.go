package auth

import (
	"test-project/internal/domain/user"
	"time"
)

type Service interface {
	Register(email, password string) (user.User, error)
	Login(email, password string) (token string, err error)
	// помечаем пользователя онлайн (вызыв. на каждом запросе)
	TouchOnline(userID string) error
	// возвращает список онлайн-ID
	OnlineUsers(since time.Duration) ([]string, error)
}
