package usecase

import (
	"errors"
	"strings"
	userDomain "test-project/internal/domain/user"
	"test-project/internal/validator"
)

type UserUsecase interface {
	ListUsers() ([]userDomain.User, error)
	GetUser(id string) (userDomain.User, error)
	DeleteUser(id string) error
	CreateUser(input userDomain.User) (userDomain.User, error)
}

type userUsecase struct {
	repo      userDomain.UserRepository
	validator *validator.Validator
}

func NewUserUsecase(r userDomain.UserRepository, v *validator.Validator) UserUsecase {
	return &userUsecase{repo: r, validator: v}
}

func (u *userUsecase) ListUsers() ([]userDomain.User, error) {
	return u.repo.FindAll()
}

func (u *userUsecase) GetUser(id string) (userDomain.User, error) {
	return u.repo.FindByID(id)
}

func (u *userUsecase) CreateUser(input userDomain.User) (userDomain.User, error) {
	if errs := u.validator.Validate(input); len(errs) > 0 {
		return userDomain.User{}, errors.New(strings.Join(errs, "; "))
	}

	return u.repo.Create(input)
}

func (u *userUsecase) DeleteUser(id string) error {
	_, err := u.repo.FindByID(id)

	if err != nil {
		return err
	}

	return u.repo.Delete(id)
}
