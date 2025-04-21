package usecase

import "test-project/internal/domain"

type UserUsecase interface {
	ListUsers() ([]domain.User, error)
	GetUser(id string) (domain.User, error)
	CreateUser(input domain.User) (domain.User, error)
}

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) ListUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecase) GetUser(id string) (domain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecase) CreateUser(input domain.User) (domain.User, error) {
	// валидировать input…
	return u.repo.Create(input)
}
