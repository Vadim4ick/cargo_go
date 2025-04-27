package invitation

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInvitationRepo struct {
	db *pgxpool.Pool
}

func NewPostgresCargoRepo(db *pgxpool.Pool) InvitationRepository {
	return &PostgresInvitationRepo{db: db}
}

func (r *PostgresInvitationRepo) Create(i Invitation) (Invitation, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO invitations (email) VALUES ($1) RETURNING id", i.Email).Scan(&i.ID)

	if err != nil {
		return Invitation{}, err
	}

	return i, nil
}
