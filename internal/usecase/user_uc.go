package usecase

import (
	"errors"
	"strings"
	"test-project/internal/domain"
	"test-project/internal/validator"
)

type UserUsecase interface {
	ListUsers() ([]domain.User, error)
	GetUser(id string) (domain.User, error)
	CreateUser(input domain.User) (domain.User, error)
}

type userUsecase struct {
	repo      domain.UserRepository
	validator *validator.Validator
}

func NewUserUsecase(r domain.UserRepository, v *validator.Validator) UserUsecase {
	return &userUsecase{repo: r, validator: v}
}

func (u *userUsecase) ListUsers() ([]domain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecase) GetUser(id string) (domain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecase) CreateUser(input domain.User) (domain.User, error) {
	if errs := u.validator.Validate(input); len(errs) > 0 {
		return domain.User{}, errors.New(strings.Join(errs, "; "))
	}

	return u.repo.Create(input)
}
