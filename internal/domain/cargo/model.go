package cargo

import "time"

type Cargo struct {
	ID int `json:"id"`

	CargoNumber        string     `json:"cargoNumber" validate:"required"`
	Date               *time.Time `json:"date,omitempty"`
	LoadUnloadDate     *time.Time `json:"loadUnloadDate,omitempty"`
	Driver             string     `json:"driver" validate:"required"`
	TransportationInfo string     `json:"transportationInfo" validate:"required"`
	PayoutAmount       *float64   `json:"payoutAmount,omitempty"`
	PayoutDate         *time.Time `json:"payoutDate,omitempty"`
	PaymentStatus      *string    `json:"paymentStatus,omitempty"`
	PayoutTerms        *string    `json:"payoutTerms,omitempty"`

	CreatedAt time.Time `json:"createdAt"`

	TruckID int `json:"truckId" validate:"required"`
}

type CargoRepository interface {
	Create(cargo Cargo) (Cargo, error)
	FindAll() ([]Cargo, error)
	FindByID(id int) (Cargo, error)
}

type CreateResponse struct {
	Message string `json:"message" example:"Груз успешно создан"`
	Data    Cargo  `json:"data"`
}

type ListResponse struct {
	Message string  `json:"message" example:"Список всех грузов"`
	Data    []Cargo `json:"data"`
}

type GetResponse struct {
	Message string `json:"message" example:"Данные о грузе"`
	Data    Cargo  `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Невалидный формат JSON"`
	Data    interface{} `json:"data"`
}
