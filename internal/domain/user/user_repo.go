package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresUserRepo struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepo(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) FindAll() ([]User, error) {
	rows, err := r.db.Query(context.Background(), `SELECT id, username, email, password, role, "createdAt" FROM users`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepo) FindByID(id string) (User, error) {
	var u User
	err := r.db.QueryRow(
		context.Background(),
		`SELECT id, username, email, password, role, "createdAt" 
		   FROM users 
		  WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("Пользователь с id=%d не существует", id)
		}

		return User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepo) Create(u User) (User, error) {
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password) 
		 VALUES ($1, $2, $3)
		 RETURNING id, role, "createdAt"`,
		u.Username, u.Email, u.Password,
	).Scan(&u.ID, &u.Role, &u.CreatedAt)

	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepo) FindByEmail(email string) (User, error) {
	var u User

	const query = `
        SELECT id,
               username,
               email,
               password,
               role,
               "createdAt"
          FROM users
         WHERE email = $1
    `

	err := r.db.QueryRow(context.Background(), query, email).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.Role,
			&u.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("Пользователь с email=%s не существует", email)
		}
		return User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepo) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	return err
}
