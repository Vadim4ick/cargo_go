package usecase

import (
	"errors"
	"strings"
	"test-project/internal/domain/cargo"
	truckDomain "test-project/internal/domain/truck"
	"test-project/internal/validator"
)

type TruckUsecase interface {
	CreateTruck(input truckDomain.Truck) (truckDomain.Truck, error)
	ListTrucks() ([]truckDomain.Truck, error)
	GetTruck(id string) (truckDomain.Truck, error)
	GetTruckCargos(id string, limit int, page int) ([]cargo.Cargo, error)
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

func (u *truckUsecase) ListTrucks() ([]truckDomain.Truck, error) {
	return u.repo.FindAll()
}

func (u *truckUsecase) GetTruck(id string) (truckDomain.Truck, error) {
	return u.repo.FindByID(id)
}

func (u *truckUsecase) GetTruckCargos(id string, limit int, page int) ([]cargo.Cargo, error) {
	return u.repo.GetTruckCargos(id, limit, page)
}
