package truck

import "time"

type Truck struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type TruckRepository interface {
	Create(truck Truck) (Truck, error)
}

type CreateResponse struct {
	Message string `json:"message" example:"Машина успешно создана"`
	Data    Truck  `json:"data"`
}
