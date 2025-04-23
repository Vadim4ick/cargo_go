package domain

import "time"

type Truck struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type TruckRepository interface {
	Create(truck Truck) (Truck, error)
}
