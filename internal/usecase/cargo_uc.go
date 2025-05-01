package usecase

import (
	"errors"
	"fmt"
	"strings"
	cargoDomain "test-project/internal/domain/cargo"
	"test-project/internal/validator"
)

type CargoUsecase interface {
	CreateCargo(input cargoDomain.Cargo) (cargoDomain.Cargo, error)
	PatchCargo(input cargoDomain.UpdateCargoInput, id string) (cargoDomain.Cargo, error)
	ListGargos() ([]cargoDomain.Cargo, error)
	DeleteCargo(id string) error
	GetCargo(id string) (cargoDomain.Cargo, error)
}

type cargoUsecase struct {
	repo      cargoDomain.CargoRepository
	validator *validator.Validator
}

func NewCargoUsecase(r cargoDomain.CargoRepository, v *validator.Validator) CargoUsecase {
	return &cargoUsecase{repo: r, validator: v}
}

func (u *cargoUsecase) CreateCargo(input cargoDomain.Cargo) (cargoDomain.Cargo, error) {
	if errs := u.validator.Validate(input); len(errs) > 0 {
		return cargoDomain.Cargo{}, errors.New(strings.Join(errs, "; "))
	}

	return u.repo.Create(input)
}

func (u *cargoUsecase) ListGargos() ([]cargoDomain.Cargo, error) {
	return u.repo.FindAll()
}

func (u *cargoUsecase) GetCargo(id string) (cargoDomain.Cargo, error) {
	cargo, err := u.repo.FindByID(id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return cargo, err
}

func (u *cargoUsecase) PatchCargo(input cargoDomain.UpdateCargoInput, id string) (cargoDomain.Cargo, error) {
	return u.repo.Update(input, id)
}

func (u *cargoUsecase) DeleteCargo(id string) error {
	_, err := u.repo.FindByID(id)

	if err != nil {
		return err
	}

	return u.repo.Delete(id)
}
