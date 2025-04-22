package repository

import (
	"context"
	"errors"
	"test-project/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepo struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepo(db *pgxpool.Pool) domain.UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) FindAll() ([]domain.User, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, username, email, password, role, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepo) FindByID(id string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(context.Background(),
		"SELECT id, username, email, password, role, created_at FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepo) Create(u domain.User) (domain.User, error) {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password) 
		 VALUES ($1, $2, $3)
		 RETURNING id, role, "createdAt"`,
		u.Username, u.Email, u.Password,
	).Scan(&u.ID, &u.Role, &u.CreatedAt)

	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
