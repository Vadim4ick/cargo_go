package truck

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTruckRepo struct {
	db *pgxpool.Pool
}

func NewPostgresTruckRepo(db *pgxpool.Pool) TruckRepository {
	return &PostgresTruckRepo{db: db}
}

func (r *PostgresTruckRepo) Create(u Truck) (Truck, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO trucks (name) VALUES ($1) RETURNING id", u.Name).Scan(&u.ID)

	if err != nil {
		return Truck{}, err
	}

	return u, nil
}
