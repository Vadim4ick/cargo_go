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
	err := r.db.QueryRow(context.Background(), "INSERT INTO invitations (email, token) VALUES ($1, $2) RETURNING id", i.Email, i.Token).Scan(&i.ID)

	if err != nil {
		return Invitation{}, err
	}

	return i, nil
}
