package usecase

import (
	"errors"
	"strings"
	"test-project/internal/domain"
	"test-project/internal/validator"
)

type TruckUsecase interface {
	CreateTruck(input domain.Truck) (domain.Truck, error)
}

type truckUsecase struct {
	repo      domain.TruckRepository
	validator *validator.Validator
}

func NewTruckUsecase(r domain.TruckRepository, v *validator.Validator) TruckUsecase {
	return &truckUsecase{repo: r, validator: v}
}

func (u *truckUsecase) CreateTruck(input domain.Truck) (domain.Truck, error) {
	if errs := u.validator.Validate(input); len(errs) > 0 {
		return domain.Truck{}, errors.New(strings.Join(errs, "; "))
	}

	return u.repo.Create(input)
}
