package truck

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *PostgresTruckRepo) FindAll() ([]Truck, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name FROM trucks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trucks []Truck
	for rows.Next() {
		var t Truck
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		trucks = append(trucks, t)
	}

	return trucks, nil
}

func (r *PostgresTruckRepo) FindByID(id string) (Truck, error) {
	var t Truck
	err := r.db.QueryRow(context.Background(), "SELECT id, name FROM trucks WHERE id = $1", id).Scan(&t.ID, &t.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Truck{}, fmt.Errorf("Машина с id=%d не существует", id)
		}

		return Truck{}, err
	}

	return t, err
}
