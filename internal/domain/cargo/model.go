package cargo

import "time"

type Cargo struct {
	ID string `json:"id"`

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

	TruckID string `json:"truckId" validate:"required"`
}

type UpdateCargoInput struct {
	CargoNumber        *string    `json:"cargoNumber,omitempty"`
	Date               *time.Time `json:"date,omitempty"`
	LoadUnloadDate     *time.Time `json:"loadUnloadDate,omitempty"`
	Driver             *string    `json:"driver,omitempty"`
	TransportationInfo *string    `json:"transportationInfo,omitempty"`
	PayoutAmount       *float64   `json:"payoutAmount,omitempty"`
	PayoutDate         *time.Time `json:"payoutDate,omitempty"`
	PaymentStatus      *string    `json:"paymentStatus,omitempty"`
	PayoutTerms        *string    `json:"payoutTerms,omitempty"`
	TruckID            *string    `json:"truckId,omitempty"`
}

type CargoRepository interface {
	Create(cargo Cargo) (Cargo, error)
	FindAll() ([]Cargo, error)
	FindByID(id string) (Cargo, error)
	Update(cargo UpdateCargoInput, id string) (Cargo, error)
	Delete(id string) error
}

type CreateRequest struct {
	CargoNumber        string     `json:"cargoNumber" validate:"required" example:"1234"`
	Driver             string     `json:"driver" validate:"required" example:"Иванов Иван Иванович"`
	TransportationInfo string     `json:"transportationInfo" validate:"required" example:"Грузоперевозка"`
	TruckID            string     `json:"truckId" validate:"required" example:"1"`
	Date               *time.Time `json:"date,omitempty" example:"2023-01-01T00:00:00Z"`
	LoadUnloadDate     *time.Time `json:"loadUnloadDate,omitempty" example:"2023-01-01T00:00:00Z"`
	PayoutAmount       *float64   `json:"payoutAmount,omitempty" example:"1000"`
	PayoutDate         *time.Time `json:"payoutDate,omitempty" example:"2023-01-01T00:00:00Z"`
	PaymentStatus      *string    `json:"paymentStatus,omitempty" example:"paid"`
	PayoutTerms        *string    `json:"payoutTerms,omitempty" example:"cash"`
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

type DeleteResponse struct {
	Message string      `json:"message" example:"Груз успешно удалён"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string      `json:"message" example:"Невалидный формат JSON"`
	Data    interface{} `json:"data"`
}
