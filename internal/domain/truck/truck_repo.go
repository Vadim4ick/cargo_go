package truck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"test-project/internal/domain/cargo"

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

func (r *PostgresTruckRepo) GetTruckCargos(
	id string, limit, page int,
) ([]cargo.Cargo, error) {

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

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

	/* -- агрегируем все файлы, привязанные к cargo -- */
	COALESCE(
	  json_agg(
	    json_build_object(
	      'id',         f.id,
	      'url',        f.url,
	      'cargoId',    f.owner_id,
	      'createdAt',  f.created_at
	    )
	    ORDER BY f.created_at
	  ) FILTER (WHERE f.id IS NOT NULL),
	  '[]'
	) AS photos_json

FROM cargos c
LEFT JOIN files f
  ON     f.owner_table = 'cargos'
     AND f.owner_id    = c.id

WHERE c.truckid = $1
GROUP BY c.id
ORDER BY c."createdAt" DESC
LIMIT  $2
OFFSET $3;
`

	rows, err := r.db.Query(context.Background(), q, id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cargos []cargo.Cargo
	for rows.Next() {
		var (
			c          cargo.Cargo
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
			&photosJSON, // JSON-массив из запроса
		); err != nil {
			return nil, err
		}

		// распаковываем JSON в срез CargoPhoto
		if err := json.Unmarshal(photosJSON, &c.CargoPhotos); err != nil {
			return nil, err
		}

		cargos = append(cargos, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return cargos, nil
}
