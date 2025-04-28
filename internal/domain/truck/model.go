package truck

import (
	"test-project/internal/domain/cargo"
	"time"
)

type Truck struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type TruckRepository interface {
	Create(truck Truck) (Truck, error)
	FindAll() ([]Truck, error)
	FindByID(id string) (Truck, error)
	GetTruckCargos(id string, limit int, page int) ([]cargo.Cargo, error)
}

type CreateRequest struct {
	Name string `json:"name" example:"Машина"`
}

type CreateResponse struct {
	Message string `json:"message" example:"Машина успешно создана"`
	Data    Truck  `json:"data"`
}

type ListResponse struct {
	Message string  `json:"message" example:"Список всех машин"`
	Data    []Truck `json:"data"`
}

type GetResponse struct {
	Message string  `json:"message" example:"Машина"`
	Data    []Truck `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Невалидный формат JSON"`
	Data    interface{} `json:"data"`
}
