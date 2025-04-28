package usecase

import (
	"errors"
	"time"

	userDomain "test-project/internal/domain/user"
	"test-project/internal/redis"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password string) (userDomain.User, error)
	Login(email, password string) (string, string, error)
	TouchOnline(userID string) error
	OnlineUsers(since time.Duration) ([]string, error)

	GetUser(id string) (userDomain.User, error)
	FindByEmail(email string) (userDomain.User, error)
}

type usecase struct {
	repo  userDomain.UserRepository
	jwt   *JwtUsecase
	redis *redis.Client
}

func NewService(r userDomain.UserRepository, j *JwtUsecase, rc *redis.Client) AuthUsecase {
	return &usecase{repo: r, jwt: j, redis: rc}
}

func (u *usecase) Register(email, password string) (userDomain.User, error) {
	if _, err := u.repo.FindByEmail(email); err == nil {
		return userDomain.User{}, errors.New("Пользователь с таким email уже существует")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return u.repo.Create(
		userDomain.User{
			Email:    email,
			Password: string(hash),
		},
	)
}

func (u *usecase) Login(email, password string) (string, string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", "", errors.New("Неверный пароль")
	}

	accessToken, err := u.jwt.GenerateAccess(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.jwt.GenerateRefresh(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *usecase) FindByEmail(email string) (userDomain.User, error) {
	return u.repo.FindByEmail(email)
}

func (u *usecase) TouchOnline(userID string) error {
	// сохраняем метку времени в Redis key: "online:<userID>"
	return u.redis.SetEX("online:"+userID, "1", 5*time.Minute) // TTL 5m
}

func (u *usecase) OnlineUsers(since time.Duration) ([]string, error) {
	// простая реализация: сканируем ключи "online:*"
	return u.redis.Keys("online:*")
}

func (s *usecase) GetUser(id string) (userDomain.User, error) {
	return s.repo.FindByID(id)
}
