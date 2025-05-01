package file

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgRepo struct{ db *pgxpool.Pool }

func NewRepo(db *pgxpool.Pool) Repository { return &pgRepo{db} }

func (r *pgRepo) Create(ctx context.Context, rec Record) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO files(id, owner_id, owner_table, url)
		 VALUES ($1,$2,$3,$4)`,
		rec.ID, rec.OwnerID, rec.OwnerTable, rec.URL)
	return err
}

func (r *pgRepo) DeleteByIDs(ctx context.Context, ids []string) ([]Record, error) {
	rows, err := r.db.Query(ctx,
		`DELETE FROM files WHERE id = ANY($1) RETURNING id, owner_id, owner_table, url`,
		ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.OwnerID, &rec.OwnerTable, &rec.URL); err != nil {
			return nil, err
		}
		res = append(res, rec)
	}
	return res, nil
}

func (r *pgRepo) GetByOwner(ctx context.Context, table, ownerID string) ([]Record, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, url FROM files WHERE owner_table=$1 AND owner_id=$2`,
		table, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.URL); err != nil {
			return nil, err
		}
		list = append(list, rec)
	}
	return list, nil
}
