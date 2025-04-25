package usecase

import (
	"errors"
	"strings"
	truckDomain "test-project/internal/domain/truck"
	"test-project/internal/validator"
)

type TruckUsecase interface {
	CreateTruck(input truckDomain.Truck) (truckDomain.Truck, error)
}

type truckUsecase struct {
	repo      truckDomain.TruckRepository
	validator *validator.Validator
}

func NewTruckUsecase(r truckDomain.TruckRepository, v *validator.Validator) TruckUsecase {
	return &truckUsecase{repo: r, validator: v}
}

func (u *truckUsecase) CreateTruck(input truckDomain.Truck) (truckDomain.Truck, error) {
	if errs := u.validator.Validate(input); len(errs) > 0 {
		return truckDomain.Truck{}, errors.New(strings.Join(errs, "; "))
	}

	return u.repo.Create(input)
}
