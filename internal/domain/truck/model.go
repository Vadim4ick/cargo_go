package truck

import "time"

type Truck struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type TruckRepository interface {
	Create(truck Truck) (Truck, error)
	FindAll() ([]Truck, error)
	FindByID(id int) (Truck, error)
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
