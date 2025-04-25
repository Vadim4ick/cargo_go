package cargo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCargoRepo struct {
	db *pgxpool.Pool
}

func NewPostgresCargoRepo(db *pgxpool.Pool) CargoRepository {
	return &PostgresCargoRepo{db: db}
}

func (r *PostgresCargoRepo) Create(c Cargo) (Cargo, error) {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO cargos 
	(cargoNumber, date, loadUnloadDate, driver, transportationInfo, payoutAmount, payoutDate, paymentStatus, payoutTerms, truckId) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING id, cargoNumber, date, loadUnloadDate, driver, transportationInfo, payoutAmount, payoutDate, paymentStatus, payoutTerms, "createdAt", truckId`,
		c.CargoNumber, c.Date, c.LoadUnloadDate, c.Driver, c.TransportationInfo, c.PayoutAmount, c.PayoutDate, c.PaymentStatus, c.PayoutTerms, c.TruckID,
	).Scan(
		&c.ID,
		&c.CargoNumber,
		&c.Date,
		&c.LoadUnloadDate,
		&c.Driver,
		&c.TransportationInfo,
		&c.PayoutAmount,
		&c.PayoutDate,
		&c.PaymentStatus,
		&c.PayoutTerms,
		&c.CreatedAt,
		&c.TruckID,
	)

	if err != nil {
		return Cargo{}, err
	}

	return c, nil
}

func (r *PostgresCargoRepo) FindAll() ([]Cargo, error) {
	rows, err := r.db.Query(context.Background(), `SELECT * FROM cargos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cargos []Cargo
	for rows.Next() {
		var c Cargo
		if err := rows.Scan(
			&c.ID,
			&c.CargoNumber,
			&c.Date,
			&c.LoadUnloadDate,
			&c.Driver,
			&c.TransportationInfo,
			&c.PayoutAmount,
			&c.PayoutDate,
			&c.PaymentStatus,
			&c.PayoutTerms,
			&c.CreatedAt,
			&c.TruckID,
		); err != nil {
			return nil, err
		}
		cargos = append(cargos, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cargos, nil
}

func (r *PostgresCargoRepo) FindByID(id int) (Cargo, error) {
	var c Cargo
	err := r.db.QueryRow(context.Background(), `SELECT * FROM cargos WHERE id=$1`, id).Scan(
		&c.ID,
		&c.CargoNumber,
		&c.Date,
		&c.LoadUnloadDate,
		&c.Driver,
		&c.TransportationInfo,
		&c.PayoutAmount,
		&c.PayoutDate,
		&c.PaymentStatus,
		&c.PayoutTerms,
		&c.CreatedAt,
		&c.TruckID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Cargo{}, fmt.Errorf("груз с id=%d не существует", id)
		}
		return Cargo{}, err
	}

	return c, nil
}
