package cargo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	const query = `
SELECT
    c.id,
    c.cargonumber          AS "cargoNumber",
    c.date,
    c.loadunloaddate       AS "loadUnloadDate",
    c.driver,
    c.transportationinfo   AS "transportationInfo",
    c.payoutamount         AS "payoutAmount",
    c.payoutdate           AS "payoutDate",
    c.paymentstatus        AS "paymentStatus",
    c.payoutterms          AS "payoutTerms",
    c."createdAt",
    c.truckid              AS "truckId",
    COALESCE(
      json_agg(
        json_build_object(
          'id',        f.id,
          'url',       f.url,
          'cargoId',   f.owner_id,
          'createdAt', f.created_at
        )
        ORDER BY f.created_at
      ) FILTER (WHERE f.id IS NOT NULL),
      '[]'
    ) AS photos_json
FROM   cargos c
LEFT   JOIN files f
       ON f.owner_table = 'cargos'
      AND f.owner_id    = c.id
GROUP  BY c.id
ORDER  BY c."createdAt" DESC;
`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query cargos: %w", err)
	}
	defer rows.Close()

	var result []Cargo

	for rows.Next() {
		var (
			c          Cargo
			photosJSON []byte
		)

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
			&photosJSON,
		); err != nil {
			return nil, fmt.Errorf("scan cargo: %w", err)
		}

		// распаковываем JSON-массив файлов
		if err := json.Unmarshal(photosJSON, &c.CargoPhotos); err != nil {
			return nil, fmt.Errorf("unmarshal photos: %w", err)
		}

		result = append(result, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return result, nil
}

func (r *PostgresCargoRepo) FindByID(id string) (Cargo, error) {
	const q = `
SELECT
    c.id,
    c.cargonumber          AS "cargoNumber",
    c.date,
    c.loadunloaddate       AS "loadUnloadDate",
    c.driver,
    c.transportationinfo   AS "transportationInfo",
    c.payoutamount         AS "payoutAmount",
    c.payoutdate           AS "payoutDate",
    c.paymentstatus        AS "paymentStatus",
    c.payoutterms          AS "payoutTerms",
    c."createdAt",
    c.truckid              AS "truckId",
    COALESCE(
      json_agg(
        json_build_object(
          'id',        f.id,
          'url',       f.url,
          'cargoId',   f.owner_id,
          'createdAt', f.created_at
        )
        ORDER BY f.created_at
      ) FILTER (WHERE f.id IS NOT NULL),
      '[]'
    ) AS photos_json
FROM   cargos c
LEFT   JOIN files f
       ON f.owner_table = 'cargos'
      AND f.owner_id    = c.id
WHERE  c.id = $1
GROUP  BY c.id;`

	var (
		c          Cargo
		photosJSON []byte
	)

	err := r.db.QueryRow(context.Background(), q, id).Scan(
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
		&photosJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Cargo{}, fmt.Errorf("cargo with id=%s not found", id)
		}
		return Cargo{}, err
	}

	if err := json.Unmarshal(photosJSON, &c.CargoPhotos); err != nil {
		return Cargo{}, fmt.Errorf("decode photos: %w", err)
	}

	return c, nil
}

func (r *PostgresCargoRepo) Update(c UpdateCargoInput, id string) (Cargo, error) {
	query := "UPDATE cargos SET "
	args := []interface{}{}
	i := 1

	if c.CargoNumber != nil {
		query += fmt.Sprintf("cargoNumber = $%d, ", i)
		args = append(args, *c.CargoNumber)
		i++
	}
	if c.Date != nil {
		query += fmt.Sprintf("date = $%d, ", i)
		args = append(args, *c.Date)
		i++
	}
	if c.LoadUnloadDate != nil {
		query += fmt.Sprintf("loadUnloadDate = $%d, ", i)
		args = append(args, *c.LoadUnloadDate)
		i++
	}
	if c.Driver != nil {
		query += fmt.Sprintf("driver = $%d, ", i)
		args = append(args, *c.Driver)
		i++
	}
	if c.TransportationInfo != nil {
		query += fmt.Sprintf("transportationInfo = $%d, ", i)
		args = append(args, *c.TransportationInfo)
		i++
	}
	if c.PayoutAmount != nil {
		query += fmt.Sprintf("payoutAmount = $%d, ", i)
		args = append(args, *c.PayoutAmount)
		i++
	}
	if c.PayoutDate != nil {
		query += fmt.Sprintf("payoutDate = $%d, ", i)
		args = append(args, *c.PayoutDate)
		i++
	}
	if c.PaymentStatus != nil {
		query += fmt.Sprintf("paymentStatus = $%d, ", i)
		args = append(args, *c.PaymentStatus)
		i++
	}
	if c.PayoutTerms != nil {
		query += fmt.Sprintf("payoutTerms = $%d, ", i)
		args = append(args, *c.PayoutTerms)
		i++
	}
	if c.TruckID != nil {
		query += fmt.Sprintf("truckId = $%d, ", i)
		args = append(args, *c.TruckID)
		i++
	}

	// убрать последнюю запятую
	query = strings.TrimSuffix(query, ", ")
	// добавить WHERE
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	_, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		return Cargo{}, err
	}

	cargo, err := r.FindByID(id)
	if err != nil {
		return Cargo{}, err
	}

	return cargo, nil
}

func (r *PostgresCargoRepo) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM cargos WHERE id=$1", id)
	return err
}
