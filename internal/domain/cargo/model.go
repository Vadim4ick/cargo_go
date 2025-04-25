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
